import React, { useEffect, useState, useContext } from "react";
import logo from "../assets/images/logo.png";
import { Row, Col, Image } from "antd";
import { Typography } from "antd";
import { AppContext } from "../appcontext";
const { Text, Title } = Typography;

export var HeaderDiv = function () {
  return (
    <Row justify="space-between" align="middle" style={{ width: "100%" }}>
      <Logo />
      <MainTitle />
      <ChainState />
    </Row>
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

var ChainState = function () {
  const { chainState, setChainState } = useContext(AppContext);
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
        {chainState.block != "" ? "Latest block: " + chainState.block : ""}
        <br />
        {chainState.date !== "" ? chainState.date : ""}
        <br />
        {"Eth price: " + chainState.price}
      </Text>
    </Col>
  );
};
