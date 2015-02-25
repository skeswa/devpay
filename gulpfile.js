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

var PUBLIC_FOLDER_PATH      = path.join(__dirname, 'public'),
    FRONTEND_FOLDER_PATH    = path.join(__dirname, 'frontend'),
    BACKEND_FOLDER_PATH     = path.join(__dirname, 'backend'),

    PROJECT_NAME            = __dirname.split(path.sep).pop(),
    SERVER_EXECUTABLE       = 'server',
    DB_CONN_STRING          = 'postgres://postgres:@localhost:5432/' + PROJECT_NAME,
    SERVER_ENV              = {
        PORT: 3000,
        DB: DB_CONN_STRING,
        VERBOSE: true,
        SESSION_SECRET: 'thisisnotasecretatall'
    };

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
                    gutil.log('Failed to re-bundle client js:');
                    console.log(err);
                    if (done) done(err);
                }
            })
            .pipe(plumber())
            .pipe(source(path.join(PUBLIC_FOLDER_PATH, 'main.js')))
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
    bundler.add(path.join(FRONTEND_FOLDER_PATH, 'main.js'));
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
    bundler.add(path.join(FRONTEND_FOLDER_PATH, 'js', 'main.js'));
    // Perform initial rebundle
    return helpers.rebundle(bundler, cb);
});

// Compiles the client less
gulp.task('less', function() {
    gulp.src(path.join(FRONTEND_FOLDER_PATH, 'less', 'main.less'))
        .pipe(plumber())
        .pipe(sourcemaps.init())
        .pipe(less())
        .pipe(sourcemaps.write())
        .pipe(gulp.dest(PUBLIC_FOLDER_PATH));
});

// Condenses the pages
gulp.task('html', function() {
    gulp.src(path.join(FRONTEND_FOLDER_PATH, 'html', '**', '*.html'))
        .pipe(plumber())
        .pipe(minify({
            empty: true,
            spare: true
        }))
        .pipe(gulp.dest(PUBLIC_FOLDER_PATH, 'html'));
});

// Moves images
gulp.task('images', helpers.copyAssets('img'));
gulp.task('images-delayed', helpers.delay(helpers.copyAssets('img')));

// Clears all compiled client code
gulp.task('clean', function() {
    clean.sync(PUBLIC_FOLDER_PATH);
});

// The server process reference obj
var serverProc = undefined;

// Remove the server executable
gulp.task('clean-server', function(callback) {
    var executablePath = path.join(BACKEND_FOLDER_PATH, SERVER_EXECUTABLE);
    fs.unlinkSync(executablePath);
});

// Uses the golang toolchain to compile server sourcecode
gulp.task('compile-server', function(done) {
    var startTime = (new Date()).getTime(),
        timeDelta;

    exec('go build -o ' + SERVER_EXECUTABLE, {
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
        executablePath = path.join(BACKEND_FOLDER_PATH, SERVER_EXECUTABLE),
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
gulp.task('stop-server', function() {
    if (serverProc) {
        serverProc.kill();
    }
});

// The hollistic, atomic server task
gulp.task('server', ['clean-server', 'stop-server', 'start-server']);

// Watches changes to the client code
gulp.task('watch', ['clean', 'less', 'html', 'images', 'server', 'watchify'], function() {
    // Watch frontend stuff
    gulp.watch(path.join(FRONTEND_FOLDER_PATH, 'html', '**', '*'), ['html']);
    gulp.watch(path.join(FRONTEND_FOLDER_PATH, 'less', '**', '*'), ['less']);
    gulp.watch(path.join(FRONTEND_FOLDER_PATH, 'img', '**', '*'), ['images-delayed']);
    // Watch backend stuff
    gulp.watch(path.join(BACKEND_FOLDER_PATH, '**', '*'), ['server']);
});

// Run all compilation tasks
gulp.task('default', ['clean', 'less', 'html', 'images', 'browserify']);
