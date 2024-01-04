import React, { useState, useEffect } from "react";
import * as Wails from "../../wailsjs/runtime";
import "./inner.css";
import { Card } from "antd";
import { appearanceLine } from "./senders";

export const Recipients: React.FC = () => {
  var [asRecipient, setFromTopTen] = useState("This is the asRecipient panel");

  useEffect(() => {
    const update = (asRecipient: string) => {
      setFromTopTen(asRecipient);
    };
    Wails.EventsOn("asRecipient", update);
    return () => {
      Wails.EventsOff("asRecipient");
    };
  }, []);

  const columns = asRecipient.split(",");
  return (
    <div className="panel inner-panel-body-triple">
      <Card
        title="As Receiver"
        style={{
          textAlign: "left",
          width: "100%",
          height: "500px",
          overflowY: "auto",
          overflowX: "hidden",
        }}
      >
        {columns.map((column, index) => {
          return appearanceLine(column);
        })}
      </Card>
    </div>
  );
};
