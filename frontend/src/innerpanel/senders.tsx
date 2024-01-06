import React, { useState, useEffect, useContext } from "react";
import { AppContext } from "../appcontext";
import * as Wails from "../../wailsjs/runtime";
import "./inner.css";
import { Card } from "antd";

export const Senders: React.FC = () => {
  var [asSender, setSender] = useState("This is the asSender panel");
  var { info } = useContext(AppContext);
  var currAddr = info.split(",")[1];

  useEffect(() => {
    const update = (asSender: string) => {
      setSender(asSender);
    };
    Wails.EventsOn("asSender", update);
    return () => {
      Wails.EventsOff("asSender");
    };
  }, []);

  const columns = asSender.replace("/,$/g", "").split(",");
  return (
    <div className="panel inner-panel-body-triple">
      <Card
        title="As Sender"
        style={{
          textAlign: "left",
          width: "100%",
          height: "500px",
          overflowY: "auto",
          overflowX: "hidden",
        }}
      >
        {columns.map((column, index) => {
          return appearanceLine(column, index, currAddr);
        })}
      </Card>
    </div>
  );
};

type Tx = {
  fromA: string;
  from: string;
  toA: string;
  to: string;
  cnt: string;
};

function convert(column: string): Tx {
  var parts = column.split("|");
  return {
    fromA: parts[0] ?? "",
    from: parts[1] ?? "",
    toA: parts[2] ?? "",
    to: parts[3] ?? "",
    cnt: parts[4] ?? "",
  };
}

var shrink = function (str: string) {
  return str.substring(0, 6) + "..." + str.substring(str.length - 4);
};

export function appearanceLine(
  column: string,
  index: number,
  currAddr: string
) {
  var other = {
    fontFamily: '"Courier New", monospace',
    fontWeight: "light",
    color: "#222222" /* "#f79090", */,
  };
  var self = {
    fontFamily: '"Courier New", monospace',
    fontWeight: "bold",
    color: "darkpurple",
  };

  var tx = convert(column);
  var isFrom = tx.fromA == currAddr;
  var isTo = tx.toA == currAddr;

  var sFrom = isFrom ? self : other;
  var sTo = isTo ? self : other;

  var fC = <div>{tx.from == "" ? tx.fromA : tx.from}</div>;
  var tC = (
    <div>
      {" ==> "}
      {tx.to == "" ? tx.toA : tx.to}
    </div>
  );

  return (
    <div
      key={index}
      style={{ display: "grid", gridTemplateColumns: "8fr 1fr" }}
    >
      <div style={sFrom}>{fC}</div>
      <div>{tx.cnt}</div>
      <div style={sTo}>{tC}</div>
    </div>
  );
}
