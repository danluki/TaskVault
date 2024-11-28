import { jsx as _jsx } from "react/jsx-runtime";
import Icon from '@mui/icons-material/ThumbUp';
import CardWithIcon from './CardWithIcon';
var SuccessfulJobs = function (_a) {
    var value = _a.value;
    return (_jsx(CardWithIcon, { to: '/jobs?filter={"status":"success"}', icon: Icon, title: 'Successful Jobs', subtitle: value }));
};
export default SuccessfulJobs;
