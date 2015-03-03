import React from 'react';

import SessionStore from './stores/session';

export default {
    protect: {
        page: (AuthedComponent) => {
            let wrapper = {},
                ProtectorComponent = React.createClass({
                    render: () => {
                        if (SessionStore.isLoggedIn()) {
                            return <AuthedComponent/>;
                        } else {
                            return <wrapper.UnAuthedComponent/>;
                        }
                    }
                });

            return {
                with: (c) => {
                    wrapper.UnAuthedComponent = c;
                    // Return the renderable component
                    return ProtectorComponent;
                }
            };
        }
    }
};
