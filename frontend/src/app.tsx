import React, { useState, useEffect, useContext } from "react";
import * as Wails from "../wailsjs/runtime";
import "./App.css";
import {
  Logger,
  FooterDiv,
  HeaderDiv,
  SideBar,
  BarChart,
  MyTableComponent,
} from "./components";
import { Card, Layout, Table } from "antd";
const { Header, Footer, Sider, Content } = Layout;

export const App: React.FC = () => {
  return (
    <Layout
      style={{ backgroundColor: "rgba(33, 37, 41, 1)", minHeight: "100vh" }}
    >
      <Header
        style={{
          backgroundColor: "rgba(33, 37, 41, 1)",
          borderBottom: "1px solid lightgray",
          height: "90px",
        }}
      >
        <HeaderDiv />
      </Header>
      <Layout>
        <Sider width={200} style={{ backgroundColor: "rgba(33, 37, 41, 1)" }}>
          <SideBar />
        </Sider>
        <Content
          style={{
            backgroundColor: "rgba(33, 37, 41, 1)",
            overflow: "auto",
            height: "calc(100vh - 130px)",
          }}
        >
          <InnerPanel />
        </Content>
      </Layout>
      <Footer
        style={{
          backgroundColor: "rgba(33, 37, 41, 1)",
          borderTop: "1px solid lightgrey",
          height: "40px",
        }}
      >
        <FooterDiv />
      </Footer>
    </Layout>
  );
};

const InnerPanel = () => {
  return (
    <div className="panel inner-panel">
      <Inner />
      <Logger />
    </div>
  );
};

const Inner = function () {
  return (
    <div className="inner-panel-body">
      <Info />
      <Dates />
      <From />
      <To />
    </div>
  );
};

const Info: React.FC = () => {
  return (
    <Card className="inner-panel-body-single" style={{ width: "100%" }}>
      <MyTableComponent />;
    </Card>
  );
};

const To: React.FC = () => {
  var [toTopTen, setToTopTen] = useState("This is the toTopTen panel");
  var [toCount, setToCount] = useState("This is the toCount panel");

  useEffect(() => {
    const update = (toTopTen: string) => {
      setToTopTen(toTopTen);
    };
    Wails.EventsOn("toTopTen", update);
    return () => {
      Wails.EventsOff("toTopTen");
    };
  }, []);

  useEffect(() => {
    const update = (toCount: string) => {
      setToCount(toCount);
    };
    Wails.EventsOn("toCount", update);
    return () => {
      Wails.EventsOff("toCount");
    };
  }, []);

  const columns = toTopTen.replace("/,$/g", "").split(",");
  return (
    <div className="panel inner-panel-body-triple">
      <Card
        title="As Sender"
        style={{
          textAlign: "left",
          width: "100%",
          height: "50%",
          overflow: "auto",
        }}
      >
        {columns.map((column, index) => {
          const val = column.replace("|", "\n ==> ").replace("-", "\t");
          return <pre>{val.trim()}</pre>;
        })}
      </Card>
    </div>
  );
};

const From: React.FC = () => {
  var [fromTopTen, setFromTopTen] = useState("This is the fromTopTen panel");
  var [fromCount, setFromCount] = useState("This is the fromCount panel");

  useEffect(() => {
    const update = (fromTopTen: string) => {
      setFromTopTen(fromTopTen);
    };
    Wails.EventsOn("fromTopTen", update);
    return () => {
      Wails.EventsOff("fromTopTen");
    };
  }, []);

  useEffect(() => {
    const update = (fromCount: string) => {
      setFromCount(fromCount);
    };
    Wails.EventsOn("fromCount", update);
    return () => {
      Wails.EventsOff("fromCount");
    };
  }, []);

  const columns = fromTopTen.split(",");
  return (
    <div className="panel inner-panel-body-triple">
      <Card
        title="As Receiver"
        style={{
          textAlign: "left",
          width: "100%",
          height: "50%",
          overflow: "auto",
        }}
      >
        {columns.map((column, index) => {
          const val = column.replace("|", "\n ==> ").replace("-", "\t");
          return <pre>{val.trim()}</pre>;
        })}
      </Card>
    </div>
  );
};

const Dates: React.FC = () => {
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

interface MyCardProps {
  word: string;
  content: any;
}

const MyCard: React.FC<MyCardProps> = ({ word, content }) => {
  return (
    <div className="bordered-word-container">
      <div className="bordered-word">{word}</div>
      <div className="bordered-word-content">{content}</div>
    </div>
  );
};
