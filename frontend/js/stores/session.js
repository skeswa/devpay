import Reflux   from 'reflux';
import Actions  from '../actions';

const   LS_SESSION_DATA_KEY = 'sessionData';

let fetchDataFromStorage = () => {
    if (localStorage[LS_SESSION_DATA_KEY] !== '') {
        let sessionDataString = localStorage[LS_SESSION_DATA_KEY],
            sessionData;
        // Attempt the JSON parse
        try {
            sessionData = JSON.parse(sessionData);
            return sessionData;
        } catch(err) {
            return undefined;
        }
    } else {
        return undefined;
    }
};

let persistDataToStorage = (token, user) => {
    let sessionDataString = JSON.stringify({
        token:  token,
        user:   user
    });

    if (sessionDataString) {
        localStorage[LS_SESSION_DATA_KEY] = sessionDataString;
    }
};

let clearDataInStorage = () => {
    localStorage[LS_SESSION_DATA_KEY] = '';
}

const SessionStore = Reflux.createStore({
    listenables:    Actions,
    init:           function() {
        let sessionData = fetchDataFromStorage();

        this.token  = sessionData ? sessionData.token : undefined;
        this.user   = sessionData ? sessionData.user : undefined;
    },
    // Getters
    isLoggedIn:     function() {
        return this.token && this.user;
    },
    // Listeners
    login:          function(token, user) {
        this.token  = token;
        this.user   = user;
        // Persist new session data
        persistDataToStorage(token, user);
    },
    logout:         function() {
        this.token  = undefined;
        this.user   = undefined;
        // Persist the logout
        clearDataInStorage();
    }
});

export default SessionStore;
