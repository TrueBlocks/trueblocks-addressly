import React, { useEffect, useState } from "react";
import * as Wails from "../../wailsjs/runtime";
import logo from "../assets/images/logo.png";
import { Layout, Row, Col, Menu, Dropdown, Image } from "antd";
import { Typography } from "antd";
const { Header } = Layout;
const { Text, Title } = Typography;

export var HeaderDiv = function () {
  return (
    <Header
      style={{
        borderBottom: "1px solid lightgray",
        background: "transparent",
        padding: "5px",
        height: "90px",
        width: "100%",
        margin: "5px",
      }}
    >
      <Row justify="space-between" align="middle" style={{ width: "100%" }}>
        <Logo />
        <MainTitle />
        <ChainSelector />
      </Row>
    </Header>
  );
};

var Logo = function () {
  return (
    <Col
      flex={1}
      style={{
        textAlign: "left",
      }}
    >
      <Image src={logo} alt="logo" />
    </Col>
  );
};

var MainTitle = function () {
  return (
    <Col flex={4} style={{ textAlign: "center" }}>
      <Title level={2} style={{ color: "white" }}>
        Account Browser
      </Title>
    </Col>
  );
};

var ChainSelector = function () {
  var [price, setPrice] = useState(0.0);
  var [latest, setLatest] = useState(0);

  useEffect(() => {
    const update = (price: number) => {
      if (price > 0.0) {
        setPrice(price);
      }
    };
    Wails.EventsOn("price", update);
    return () => {
      Wails.EventsOff("price");
    };
  }, []);

  useEffect(() => {
    const update = (latest: number) => {
      if (latest > 0) {
        setLatest(latest);
      }
    };
    Wails.EventsOn("latest", update);
    return () => {
      Wails.EventsOff("latest");
    };
  }, []);

  return (
    <Col
      flex={1}
      style={{
        display: "flex", // Make Col a flex container
        flexDirection: "column", // Stack children vertically
        justifyContent: "flex-end", // Push children to the bottom
        alignItems: "flex-end", // Push children to the bottom
        paddingRight: "20px",
        height: "80px", // Set a specific height for the Col
      }}
    >
      <Text style={{ textAlign: "right", color: "white", fontSize: ".9em" }}>
        {latest > 0 ? "Latest date: " + latest : ""}
        <br />
        {latest > 0 ? "Latest block: " + latest : ""}
        <br />
        {price > 0.0 ? "Eth price: " + price : ""}
        <br />
        Selector
        {/*
        return (
          <div className="panel header-right">
            <div className="price"></div>
              <select id="chain-select">
                <option value="mainnet">Mainnet</option>
                <option value="optimism">Optimism</option>
                <option value="optimism">Sepolia</option>
                <option value="optimism">Polygon</option>
              </select>
            </div>
          </div>
        );
        */}
      </Text>
    </Col>
  );
};
