import SessionStore from '../stores/session';

export default {
    statics: {
        willTransitionTo: function(transition) {
            if (SessionStore.isLoggedIn()) {
                transition.redirect('/');
            }
        }
    }
};
