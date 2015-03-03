import React from 'react';
import {Navigation, Link} from 'react-router';
import PublicPageMixin from '../../mixins/publicpage';

const Register = React.createClass({
    mixins: [
        PublicPageMixin
    ],
    render: () => {
        return (
            <div>
                <h1>register</h1>
            </div>
        );
    }
});

export default Register;
