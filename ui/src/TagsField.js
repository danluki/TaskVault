import { jsx as _jsx } from "react/jsx-runtime";
import { Chip } from '@mui/material';
export var TagsField = function (_a) {
    var record = _a.record;
    if (record === undefined) {
        return null;
    }
    else {
        return _jsx("ul", { children: Object.keys(record.Tags).map(function (key) { return (_jsx(Chip, { label: key + ": " + record.Tags[key] })); }) });
    }
};
TagsField.defaultProps = {
    addLabel: true
};
