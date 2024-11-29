import { Card, CardContent, CardHeader } from "@mui/material";
import { List, Datagrid, TextField } from "react-admin";
import { TagsField } from "../TagsField";
import Leader from "./Leader";
// import FailedPairs from "./FailedPairs";
// import SuccessfulPairs from "./SuccessfulPairs";
// import UntriggeredPairs from "./UntriggeredPairs";
import TotalPairs from "./TotalPairs";

let fakeProps = {
  basePath: "/members",
  count: 10,
  hasCreate: false,
  hasEdit: false,
  hasList: true,
  hasShow: false,
  location: { pathname: "/", search: "", hash: "", state: undefined },
  match: { path: "/", url: "/", isExact: true, params: {} },
  options: {},
  permissions: null,
  resource: "members",
};

const styles = {
  flex: { display: "flex" },
  flexColumn: { display: "flex", flexDirection: "column" },
  leftCol: { flex: 1, marginRight: "0.5em" },
  rightCol: { flex: 1, marginLeft: "0.5em" },
  singleCol: { marginTop: "1em", marginBottom: "1em" },
};

const Spacer = () => <span style={{ width: "1em" }} />;

const Dashboard = () => (
  <div>
    <Card>
      <CardHeader title="Welcome" />
      <CardContent>
        <div style={styles.flex}>
          <div style={styles.leftCol}>
            <div style={styles.flex}>
              <Leader value={window.TASKVAULT_LEADER || "devel"} />
              <Spacer />
              <TotalPairs value={window.TASKVAULT_TOTAL_PAIRS || "0"} />
              <Spacer />
              {/* <PairsAdded value={window.TASKVAULT_PAIRS_ADDED || "0"} /> */}
              <Spacer />
              {/* <PairsUpdated value={window.TASKVAULT_PAIRS_UPDATED || "0"} /> */}
              <Spacer />
              {/* <UntriggeredPairs
                value={window.window.TASKVAULT_PAIRS_DELETED || "0"}
              /> */}
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
    <Card>
      <CardHeader title="Nodes" />
      <CardContent>
        <List {...fakeProps}>
          <Datagrid isRowSelectable={(record) => false}>
            <TextField source="Name" sortable={false} />
            <TextField source="Addr" sortable={false} />
            <TextField source="Port" sortable={false} />
            <TextField label="Status" source="statusText" sortable={false} />
            <TagsField source="Tags" sortable={false} />
          </Datagrid>
        </List>
      </CardContent>
    </Card>
  </div>
);
export default Dashboard;
