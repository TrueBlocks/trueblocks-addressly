import React, { useState, useEffect, useContext } from "react";
import { AppContext } from "../appcontext";
import * as Wails from "../../wailsjs/runtime";
import "./inner.css";
import { Card } from "antd";

export const Senders: React.FC = () => {
  var [asSender, setToTopTen] = useState("This is the asSender panel");

  useEffect(() => {
    const update = (asSender: string) => {
      setToTopTen(asSender);
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
          return appearanceLine(column);
        })}
      </Card>
    </div>
  );
};

export function appearanceLine(column: string) {
  const { address, info } = useContext(AppContext);

  var shrink = function (str: string) {
    return str.substring(0, 6) + "..." + str.substring(str.length - 4);
  };

  var red = {
    fontFamily: '"Courier New", monospace',
    fontWeight: "light",
    color: "#f0b0b0",
  };
  var blue = {
    fontFamily: '"Courier New", monospace',
    fontWeight: "bold",
    color: "dodgerblue",
  };
  var s1 = {};
  var s2 = {};
  var from, to, cnt;
  var parts = column.split("|");
  if (parts.length > 4) {
    if (parts[1] !== "") {
      from = parts[1].substring(0, 20) + " [" + shrink(parts[0]) + "]";
    } else {
      from = parts[0];
    }
    if (parts[3] !== "") {
      to = parts[3].substring(0, 20) + " [" + shrink(parts[2]) + "]";
    } else {
      to = parts[2];
    }
    cnt = parts[4];
    var iParts = info.split(",");
    if (parts[0] === iParts[1]) {
      s1 = blue;
    } else {
      s1 = red;
    }
    if (parts[2] === iParts[1]) {
      s2 = blue;
    } else {
      s2 = red;
    }
    // to = "parts[0]: " + parts[0] + " iParts[i]: " + iParts[1];
  } else {
    return <div>{column}</div>;
  }

  return (
    <div style={{ display: "grid", gridTemplateColumns: "8fr 1fr" }}>
      <div style={s1}>{from}</div>
      <div>{cnt}</div>
      <div style={s2}>
        {" ==> "}
        {to}
      </div>
      {/* <div style={{ color: "yellowgreen" }}>{info}</div> */}
    </div>
  );
  // column = from + "|" + to + "|" + cnt;
  // const val = column.replace(/\|/g, "\n ==> ");
  // return <pre>{val.trim()}</pre>;
}
