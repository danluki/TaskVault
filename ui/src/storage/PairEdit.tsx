import {
  Edit,
  SelectInput,
  TextInput,
  Create,
  SimpleForm,
  BooleanInput,
  NumberInput,
  DateTimeInput,
  required,
  useRecordContext,
} from "react-admin";
import { JsonInput } from "react-admin-json-view";

export const PairEdit = () => {
  const record = useRecordContext();
  return (
    <Edit {...record}>
      <EditForm />
    </Edit>
  );
};

export const PairCreate = (props: any) => (
  <Create {...props}>
    <EditForm />
  </Create>
);

const EditForm = (record: any) => (
  <SimpleForm {...record}>
    <TextInput source="key" helperText="Key." validate={required()} />
    <TextInput source="value" helperText="Value." validate={required()} />
  </SimpleForm>
);
