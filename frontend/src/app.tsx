import React from "react";
import "./App.css";
import { FooterDiv, HeaderDiv, SideBar } from "./components";
import { InnerPanel } from "./innerpanel";
import { Layout } from "antd";
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
