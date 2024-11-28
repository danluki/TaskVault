import { Admin, Resource, CustomRoutes } from 'react-admin';
import { Route } from "react-router-dom";
import PlayCircleOutlineIcon from '@mui/icons-material/PlayCircleOutline';
import { createHashHistory } from "history";

import { Layout } from './layout';
import dataProvider from './dataProvider';
import Dashboard from './dashboard';
import Settings from './settings/Settings';

declare global {
    interface Window {
        TASKVAULT_API_URL: string;
        TASKVAULT_LEADER: string;
        TASKVAULT_TOTAL_PAIRS: string;
        TASKVAULT_PAIRS_ADDED: string;
        TASKVAULT_PAIRS_DELETED: string;
        TASKVAULT_PAIRS_UPDATED: string;
    }
}

const history = createHashHistory();
 
export const App = () => <Admin
    dashboard={Dashboard}
    dataProvider={dataProvider}
    layout={Layout}
>

    {/* <Resource name="executions" /> */}
    <Resource name="members" />
    <CustomRoutes>
        <Route path="/settings" element={<Settings />} />
    </CustomRoutes>
</Admin>;
