import {
  Datagrid,
  TextField,
  EditButton,
  Filter,
  TextInput,
  List,
  Pagination,
} from "react-admin";
import { makeStyles } from "@mui/styles";

const PairFilter = (props: any) => (
  <Filter {...props}>
    <TextInput label="Search" source="q" alwaysOn />
    {/* <SelectInput
      source="status"
      choices={[
        { id: "success", name: "Success" },
        { id: "failed", name: "Failed" },
        { id: "untriggered", name: "Waiting to Run" },
      ]}
    /> */}
    {/* <BooleanInput source="disabled" /> */}
  </Filter>
);

const PairPagination = (props: any) => (
  <Pagination rowsPerPageOptions={[5, 10, 25, 50, 100]} {...props} />
);

const useStyles = makeStyles((theme) => ({
  hiddenOnSmallScreens: {
    display: "table-cell",
    [(theme as any).breakpoints.down("md")]: {
      display: "none",
    },
  },
  cell: {
    padding: "6px 8px 6px 8px",
  },
}));

const PairList = (props: any) => {
  const classes = useStyles();
  return (
    <List {...props} filters={<PairFilter />} pagination={<PairPagination />}>
      <Datagrid rowClick="show">
        <TextField
          source="Key"
          cellClassName={classes.hiddenOnSmallScreens}
          headerClassName={classes.hiddenOnSmallScreens}
        />
        <TextField
          source="Value"
          cellClassName={classes.hiddenOnSmallScreens}
          headerClassName={classes.hiddenOnSmallScreens}
        />
        {/* <DateField source="ttl" showTime /> */}
        <EditButton />
      </Datagrid>
    </List>
  );
};

export default PairList;
