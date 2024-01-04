import React, { useState, useEffect } from "react";
import * as Wails from "../../wailsjs/runtime";
import "./inner.css";
import { Card } from "antd";
import { appearanceLine } from "./senders";
import { resetWarned } from "antd/es/_util/warning";

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
        // title={
        //   <div style={{ display: "grid", gridTemplateColumns: "10fr 1fr" }}>
        //     <div
        //       style={{
        //         margin: 0,
        //         padding: 0,
        //         color: "blue",
        //       }}
        //     >
        //       As Received
        //     </div>
        //     <div
        //       style={{
        //         margin: 0,
        //         padding: 0,
        //         color: "blue",
        //         fontSize: ".7em",
        //         fontStyle: "italic",
        //       }}
        //     >
        //       context
        //     </div>
        //   </div>
        // }
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
