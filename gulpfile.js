var gulp        = require('gulp'),
    // Generic imports
    fs          = require('fs'),
    Stream      = require('stream'),
    gutil       = require('gulp-util'),
    path        = require('path'),
    clean       = require('rimraf'),
    plumber     = require('gulp-plumber'),
    // Browserify-related imports
    browserify  = require('browserify'),
    source      = require('vinyl-source-stream'),
    watchify    = require('watchify'),
    reactify    = require('reactify'),
    uglify      = require('gulp-uglify'),
    buffer      = require('vinyl-buffer')
    // LESS-related imports
    less        = require('gulp-less'),
    sourcemaps  = require('gulp-sourcemaps'),
    // HTML-related imports
    minify      = require('gulp-minify-html'),
    // Process-related imports
    spawn       = require('child_process').spawn,
    exec        = require('child_process').exec;


/**************************************** CONSTANTS ****************************************/

var PUBLIC_FOLDER_NAME          = 'public',
    FRONTEND_FOLDER_NAME        = 'frontend',
    FRONTEND_JS_FOLDER_NAME     = 'js',
    FRONTEND_LESS_FOLDER_NAME   = 'less',
    FRONTEND_IMG_FOLDER_NAME    = 'img',
    FRONTEND_VENDOR_FOLDER_NAME = 'vendor',
    FRONTEND_HTML_FOLDER_NAME   = 'html',
    BACKEND_FOLDER_NAME         = 'backend',

    FRONTEND_JS_ENTRY_POINT     = 'main.js',
    FRONTEND_LESS_ENTRY_POINT   = 'main.less',
    BACKEND_EXECUTABLE_NAME     = 'server',

    PUBLIC_FOLDER_PATH          = path.join(__dirname, PUBLIC_FOLDER_NAME),
    FRONTEND_FOLDER_PATH        = path.join(__dirname, FRONTEND_FOLDER_NAME),
    FRONTEND_JS_FOLDER_PATH     = path.join(FRONTEND_FOLDER_PATH, FRONTEND_JS_FOLDER_NAME),
    FRONTEND_LESS_FOLDER_PATH   = path.join(FRONTEND_FOLDER_PATH, FRONTEND_LESS_FOLDER_NAME),
    FRONTEND_IMG_FOLDER_PATH    = path.join(FRONTEND_FOLDER_PATH, FRONTEND_IMG_FOLDER_NAME),
    FRONTEND_VENDOR_FOLDER_PATH = path.join(FRONTEND_FOLDER_PATH, FRONTEND_VENDOR_FOLDER_NAME),
    FRONTEND_HTML_FOLDER_PATH   = path.join(FRONTEND_FOLDER_PATH, FRONTEND_HTML_FOLDER_NAME),
    BACKEND_FOLDER_PATH         = path.join(__dirname, BACKEND_FOLDER_NAME),

    PROJECT_NAME                = __dirname.split(path.sep).pop(),
    DB_CONN_STRING              = 'postgres://postgres:@localhost:5432/' + PROJECT_NAME,
    SERVER_ENV                  = {
        PORT:           3000,
        DB:             DB_CONN_STRING,
        VERBOSE:        true,
        SESSION_SECRET: 'thisisnotasecretatall'
    };

/************************************* HELPER FUNCTIONS ************************************/

var helpers = {
    rebundle: function(bundler, done) {
        var time = (new Date()).getTime();
        gutil.log('Started re-bundling client js');
        bundler
            .bundle(function(err) {
                if (!err) {
                    gutil.log('Finished re-bundling client js after ' + (((new Date()).getTime() - time) / 1000) + ' s');
                    if (done) done();
                } else {
                    gutil.log('Failed to re-bundle client js');
                    if (done) done(err);
                }
            })
            .pipe(plumber())
            .pipe(source(FRONTEND_JS_ENTRY_POINT))
            .pipe(buffer())
            .pipe(uglify())
            .pipe(gulp.dest(PUBLIC_FOLDER_PATH));
    },
    delay: function(callback) {
        // Waits a second before executing a function
        return function() {
            setTimeout(callback, 1000);
        };
    },
    copyAssets: function(folderName) {
        return function() {
            gulp.src(path.join(FRONTEND_FOLDER_PATH, folderName, '**', '*'))
                .pipe(plumber())
                .pipe(gulp.dest(path.join(PUBLIC_FOLDER_PATH, folderName)));
        };
    }
};

/**************************************** FRONTEND ****************************************/

// Compiles the client js
gulp.task('browserify', function(cb) {
    var bundler = browserify({
        cache: {},
        packageCache: {},
        fullPaths: true
    });
    // JSX compilation middleware
    bundler.transform(reactify);
    // Add the entry point
    bundler.add(path.join(FRONTEND_JS_FOLDER_PATH, FRONTEND_JS_ENTRY_POINT));
    // Perform initial rebundle
    return helpers.rebundle(bundler, cb);
});

// Watches and recompiles client js
gulp.task('watchify', function(cb) {
    var bundler = browserify({
        cache: {},
        packageCache: {},
        fullPaths: true,
        debug: true
    });
    // Pass the browserify bundler to watchify
    bundler = watchify(bundler);
    // JSX compilation middleware
    bundler.transform(reactify);
    // Bundlize on updates
    bundler.on('update', function() {
        helpers.rebundle(bundler);
    });
    // Add the entry point
    bundler.add(path.join(FRONTEND_JS_FOLDER_PATH, FRONTEND_JS_ENTRY_POINT));
    // Perform initial rebundle
    return helpers.rebundle(bundler, cb);
});

// Compiles the client less
gulp.task('less', function() {
    gulp.src(path.join(FRONTEND_LESS_FOLDER_PATH, FRONTEND_LESS_ENTRY_POINT))
        .pipe(plumber())
        .pipe(sourcemaps.init())
        .pipe(less())
        .pipe(sourcemaps.write())
        .pipe(gulp.dest(PUBLIC_FOLDER_PATH));
});

// Condenses the pages
gulp.task('html', function() {
    gulp.src(path.join(FRONTEND_HTML_FOLDER_PATH, '**', '*'))
        .pipe(plumber())
        .pipe(minify({
            empty: true,
            spare: true
        }))
        .pipe(gulp.dest(path.join(PUBLIC_FOLDER_PATH, FRONTEND_HTML_FOLDER_NAME)));
});

// Moves images
gulp.task('images', helpers.copyAssets(FRONTEND_IMG_FOLDER_NAME));
gulp.task('images-delayed', helpers.delay(helpers.copyAssets(FRONTEND_IMG_FOLDER_NAME)));

// Moves vendor files
gulp.task('vendor', helpers.copyAssets(FRONTEND_VENDOR_FOLDER_NAME));
gulp.task('vendor-delayed', helpers.delay(helpers.copyAssets(FRONTEND_VENDOR_FOLDER_NAME)));

// Clears all compiled client code
gulp.task('clean', function() {
    clean.sync(PUBLIC_FOLDER_PATH);
});

/**************************************** BACKEND ****************************************/

// The server process reference obj
var serverProc = undefined;

// Remove the server executable
gulp.task('clean-server', ['stop-server'], function(done) {
    var executablePath = path.join(BACKEND_FOLDER_PATH, BACKEND_EXECUTABLE_NAME);
    fs.exists(executablePath, function(exists) {
        if (exists) {
            fs.unlink(executablePath, function(err) {
                if (err) {
                    done('Could not delete the server executable');
                } else {
                    done();
                }
            });
        } else {
            done();
        }
    });
});

// Uses the golang toolchain to compile server sourcecode
gulp.task('compile-server', ['clean-server'], function(done) {
    var startTime = (new Date()).getTime(),
        timeDelta;

    exec('go build -o ' + BACKEND_EXECUTABLE_NAME, {
        cwd: BACKEND_FOLDER_PATH
    }, function(err, stdout, stderr) {
        if (err) {
            done(err);
        } else {
            timeDelta = (new Date()).getTime() - startTime;
            gutil.log('Server successfully compiled in ' + (timeDelta / 1000) + ' s');
            done();
        }
    });
});

// Runs the server executable
gulp.task('start-server', ['compile-server'], function(done) {
    var startTime = (new Date()).getTime(),
        callbackTriggered = false,
        executablePath = path.join(BACKEND_FOLDER_PATH, BACKEND_EXECUTABLE_NAME),
        timeDelta = 0;

    if (fs.existsSync(executablePath)) {
        serverProc = spawn(executablePath, [], {
            env: SERVER_ENV
        });
        // Setup listeners
        serverProc.stdout.on('data', function(data) {
            process.stdout.write(data);
        });
        serverProc.stderr.on('data', function(data) {
            process.stdout.write(data);
        });
        serverProc.on('close', function(code) {
            // Signal that the server finished running
            serverProc = undefined;
            // Report the server exit
            timeDelta = (new Date()).getTime() - startTime;
            gutil.log('Server exited with code \'' + code + '\' after ' + (timeDelta / 1000) + ' s');
            if (!callbackTriggered) {
                callbackTriggered = true;
                if (timeDelta < 100) {
                    done('The server exited suddenly (less than 100ms)');
                } else {
                    done();
                }
            }
        });
        // Finish automatically after 150ms
        setTimeout(function() {
            if (!callbackTriggered) {
                callbackTriggered = true;
                done();
            }
        }, 150);
    } else {
        done('The server has not been compiled yet');
    }
});

// Halts the server process if it exists
gulp.task('stop-server', function(done) {
    if (serverProc) {
        serverProc.kill();
        // Wait until the server process is killed
        var interval = setInterval(function() {
            if (!serverProc) {
                clearInterval(interval);
                done();
            }
        }, 50);
    } else {
        done();
    }
});

gulp.task('watch-server', ['start-server'], function(done) {
    gulp.watch(path.join(BACKEND_FOLDER_PATH, '**', '*'), ['start-server']);
});

// The hollistic, atomic server task
gulp.task('server', ['start-server']);

/**************************************** GENERIC ****************************************/

// Watches changes to the client code
gulp.task('watch', ['clean', 'less', 'html', 'images', 'vendor', 'watch-server', 'watchify'], function() {
    gulp.watch(path.join(FRONTEND_HTML_FOLDER_PATH,     '**', '*'), ['html']);
    gulp.watch(path.join(FRONTEND_LESS_FOLDER_PATH,     '**', '*'), ['less']);
    gulp.watch(path.join(FRONTEND_IMG_FOLDER_PATH,      '**', '*'), ['images-delayed']);
    gulp.watch(path.join(FRONTEND_VENDOR_FOLDER_PATH,   '**', '*'), ['vendor-delayed']);
});

// Run all compilation tasks
gulp.task('default', ['watch']);
