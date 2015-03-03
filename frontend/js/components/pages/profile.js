import React from 'react';
import {Navigation, Link} from 'react-router';
import PrivatePageMixin from '../../mixins/privatepage';

const Profile = React.createClass({
    mixins: [
        PrivatePageMixin
    ],
    render: () => {
        return (
            <div>
                <h1>profile</h1>
            </div>
        );
    }
});

export default Profile;
