import React from "react";
import "./inner.css";
import { Logging } from "../components";
import { Info, Charts, Recipients, Senders } from "./";

export const InnerPanel = () => {
  return (
    <div className="panel inner-panel">
      <div className="inner-panel-body">
        <Info />
        <Charts />
        <Recipients />
        <Senders />
      </div>
      <Logging />
    </div>
  );
};
