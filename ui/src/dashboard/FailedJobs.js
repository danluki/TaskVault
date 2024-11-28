import { jsx as _jsx } from "react/jsx-runtime";
import Icon from '@mui/icons-material/ThumbDown';
import CardWithIcon from './CardWithIcon';
var FailedJobs = function (_a) {
    var value = _a.value;
    return (_jsx(CardWithIcon, { to: '/jobs?filter={"status":"failed"}', icon: Icon, title: 'Failed Jobs', subtitle: value }));
};
export default FailedJobs;
