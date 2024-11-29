import * as React from "react";
import { FC } from "react";
import Icon from "@mui/icons-material/Update";

import CardWithIcon from "./CardWithIcon";

interface Props {
  value?: string;
}

const TotalPairs: FC<Props> = ({ value }) => {
  return (
    <CardWithIcon
      to="/storage"
      icon={Icon}
      title="Total Values"
      subtitle={value}
    />
  );
};

export default TotalPairs;
