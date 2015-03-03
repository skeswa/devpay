import React        from 'react';
import Router       from 'react-router';
// Page imports
import Login        from './components/pages/login';
import Register     from './components/pages/register';
import Profile      from './components/pages/profile';
import Splash       from './components/pages/splash';
import FourOhFour   from './components/pages/404';
import Home         from './components/pages/home';
// Misc. imports
import Util         from './util';

// React-router variables
var Route           = Router.Route,
    RouteHandler    = Router.RouteHandler,
    Redirect        = Router.Redirect,
    DefaultRoute    = Router.DefaultRoute,
    NotFoundRoute   = Router.NotFoundRoute;

// Web application page structure
let sitemap = (
    <Route handler={RouteHandler}>
        <Route name="login" handler={Login}/>
        <Route name="register" handler={Register}/>
        <Route name="profile" handler={Profile}/>
        <NotFoundRoute handler={FourOhFour}/>
        <DefaultRoute handler={Util.protect.page(Home).with(Splash)}/>
    </Route>
);

// Render to the <body/> of the DOM
Router.run(sitemap, Router.HashLocation, function(Handler) {
    React.render(<Handler/>, document.body);
});
