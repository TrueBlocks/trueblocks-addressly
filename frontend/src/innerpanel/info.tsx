import React from "react";
import "./inner.css";
import { MyTableComponent } from "../components";
import { Card } from "antd";

export const Info: React.FC = () => {
  return (
    <Card className="inner-panel-body-single" style={{ width: "100%" }}>
      <MyTableComponent />;
    </Card>
  );
};
