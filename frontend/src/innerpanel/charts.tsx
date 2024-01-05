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

  return (
    <div className="panel inner-panel-body-triple">
      <Card title="Appearances" style={{ textAlign: "left", width: "100%" }}>
        <Tabs activeKey={chartType} onChange={handleTabChange}>
          <TabPane tab="Annually" key="year">
            <BarChart str={years} />
          </TabPane>
          <TabPane tab="Monthly" key="month">
            <BarChart str={months} />
          </TabPane>
        </Tabs>
      </Card>
    </div>
  );
};
