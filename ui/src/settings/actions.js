export var CHANGE_THEME = 'CHANGE_THEME';
export var changeTheme = function (theme) { return ({
    type: CHANGE_THEME,
    payload: theme,
}); };
