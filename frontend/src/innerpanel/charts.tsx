import React, { useState, useEffect } from "react";
import * as Wails from "../../wailsjs/runtime";
import "./inner.css";
import { BarChart } from "../components";
import { Card, Tabs } from "antd";
import { SetChartType } from "../../wailsjs/go/main/App";

export const Charts: React.FC = () => {
  var [years, setYears] = useState("This is the years panel");
  var [months, setMonths] = useState("This is the months panel");
  var [chartType, setChartType] = useState("");

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

  useEffect(() => {
    const update = (chartType: string) => {
      setChartType(chartType);
    };
    Wails.EventsOn("chartType", update);
    return () => {
      Wails.EventsOff("chartType");
    };
  }, []);

  const handleTabChange = (chartType: string) => {
    setChartType(chartType);
    SetChartType(chartType);
  };

  const tabItems = [
    {
      label: "Annually",
      key: "year",
      children: <BarChart str={years} />,
    },
    {
      label: "Monthly",
      key: "month",
      children: <BarChart str={months} />,
    },
  ];

  return (
    <div className="panel inner-panel-body-triple">
      <Card title="Frequency" style={{ textAlign: "left", width: "100%" }}>
        <Tabs
          activeKey={chartType}
          onChange={handleTabChange}
          items={tabItems}
        />
      </Card>
    </div>
  );
};
