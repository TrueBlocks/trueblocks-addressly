import React, { useState, useEffect } from "react";
import * as Wails from "../../wailsjs/runtime";
import "./inner.css";
import { BarChart } from "../components";
import { Card } from "antd";

export const Charts: React.FC = () => {
  var [years, setYears] = useState("This is the years panel");
  var [months, setMonths] = useState("This is the months panel");

  useEffect(() => {
    const update = (months: string) => {
      setMonths(months);
    };
    Wails.EventsOn("months", update);
    return () => {
      Wails.EventsOff("months");
    };
  }, []);

  useEffect(() => {
    const update = (years: string) => {
      setYears(years);
    };
    Wails.EventsOn("years", update);
    return () => {
      Wails.EventsOff("years");
    };
  }, []);

  return (
    <div className="panel inner-panel-body-triple">
      <Card
        title="Appearances per Year"
        style={{
          textAlign: "left",
          width: "100%",
          marginBottom: "10px",
        }}
      >
        <BarChart str={years} />
      </Card>
      <Card
        title="Appearances per Month"
        style={{
          textAlign: "left",
          width: "100%",
        }}
      >
        <BarChart str={months} />
      </Card>
    </div>
  );
};
