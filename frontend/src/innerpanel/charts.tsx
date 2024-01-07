import React, { useState, useEffect, useContext } from "react";
import * as Wails from "../../wailsjs/runtime";
import "./inner.css";
import { AppContext } from "../appcontext";
import { BarChart } from "../components";
import { Card, Tabs } from "antd";
import { SetChartType } from "../../wailsjs/go/main/App";

export const Charts: React.FC = () => {
  const [years, setYears] = useState("This is the years panel");
  const [months, setMonths] = useState("This is the months panel");
  const { chartType, setChartType } = useContext(AppContext);

  useEffect(() => {
    const update = (months: string): void => {
      setMonths(months);
    };
    Wails.EventsOn("months", update);
    return () => {
      Wails.EventsOff("months");
    };
  }, []);

  useEffect(() => {
    const update = (years: string): void => {
      setYears(years);
    };
    Wails.EventsOn("years", update);
    return () => {
      Wails.EventsOff("years");
    };
  }, []);

  const handleTabChange = (ct: string): void => {
    setChartType(ct);
    SetChartType(ct);
  };

  const tabItems = [
    {
      label: "Annually",
      key: "year",
      children: (
        <>
          <div>{years}</div>
          <BarChart dataStr={years} clickHandler={handleTabChange} />
        </>
      )
    },
    {
      label: "Monthly",
      key: "month",
      children: (
        <>
          <div>{months}</div>
          <BarChart dataStr={months} clickHandler={handleTabChange} />
        </>
      )
    }
  ];

  return (
    <div className="panel inner-panel-body-triple">
      <Card
        title="Frequency"
        style={{
          textAlign: "left",
          width: "100%"
        }}
      >
        <Tabs
          activeKey={chartType}
          onChange={handleTabChange}
          items={tabItems}
        />
      </Card>
    </div>
  );
};
