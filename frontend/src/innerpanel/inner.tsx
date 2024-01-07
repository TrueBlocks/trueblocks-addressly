import React, { useEffect, type ReactElement } from "react";
import * as Wails from "../../wailsjs/runtime";
import "./inner.css";
import { Logging } from "../components";
import { Info, Charts, Recipients, Senders } from "./";
import { Tabs, Card } from "antd";

export const InnerPanel = (): ReactElement => {
  const [txType, setTxType] = React.useState("recipients");

  useEffect(() => {
    const update = (t: string): void => {
      setTxType(t);
    };
    Wails.EventsOn("txType", update);
    return () => {
      Wails.EventsOff("txType");
    };
  }, []);

  const handleTabChange = (t: string): void => {
    setTxType(t);
    // SetTxType(t);
    // console.log("handleTabChange txType: ", t, txType);
  };

  const tabItems = [
    {
      label: "as Recipient",
      key: "recipients",
      children: <Recipients />
    },
    {
      label: "as Sender",
      key: "senders",
      children: <Senders />
    }
  ];

  return (
    <div className="panel inner-panel">
      <div className="inner-panel-body">
        <Info />
        <Charts />
        <Card
          title="Transaction Counts"
          style={{
            textAlign: "left",
            width: "100%",
            overflowY: "auto",
            overflowX: "hidden",
            padding: "10px"
          }}
        >
          <Tabs
            activeKey={txType}
            onChange={handleTabChange}
            items={tabItems}
          />
        </Card>
      </div>
      <Logging />
    </div>
  );
};
