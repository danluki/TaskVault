import { jsx as _jsx } from "react/jsx-runtime";
import Icon from '@mui/icons-material/NewReleases';
import CardWithIcon from './CardWithIcon';
var UntriggeredJobs = function (_a) {
    var value = _a.value;
    return (_jsx(CardWithIcon, { to: '/jobs?filter={"status":"untriggered"}', icon: Icon, title: 'Untriggered Jobs', subtitle: value }));
};
export default UntriggeredJobs;
