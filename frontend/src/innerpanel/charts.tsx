import React, { useState, useEffect } from "react";
import * as Wails from "../../wailsjs/runtime";
import "./inner.css";
import { BarChart } from "../components";
import { Card, Tabs } from "antd";
import { SetChartType } from "../../wailsjs/go/main/App";

const { TabPane } = Tabs;

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

  const handleTabChange = (key: string) => {
    const chartType = key === "1" ? "Year" : "Month";
    SetChartType(chartType);
  };

  return (
    <div className="panel inner-panel-body-triple">
      <Card title="Appearances" style={{ textAlign: "left", width: "100%" }}>
        <Tabs defaultActiveKey="1" onChange={handleTabChange}>
          <TabPane tab="Annually" key="1">
            <BarChart str={years} />
          </TabPane>
          <TabPane tab="Monthly" key="2">
            <BarChart str={months} />
          </TabPane>
        </Tabs>
      </Card>
    </div>
  );
};
