import React from 'react';
import {Navigation, Link} from 'react-router';
import PublicPageMixin from '../../mixins/publicpage';

const Login = React.createClass({
    mixins: [
        PublicPageMixin
    ],
    render: () => {
        return (
            <div>
                <h1>login</h1>
            </div>
        );
    }
});

export default Login;
