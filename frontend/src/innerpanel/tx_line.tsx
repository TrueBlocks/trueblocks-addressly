import React from "react";
import { Table } from "antd";
import { ArrowRightOutlined } from "@ant-design/icons";
import "./inner.css";

export const TransactionsTable: React.FC<{
  data: string[];
  currAddr: string;
}> = ({ data, currAddr }) => {
  const dataSource = data.map((item, index) => ({
    key: index,
    transaction: str2Tx(item, currAddr)
  }));

  var shorten = function (str: string) {
    var wid = 10;
    return str.substring(0, wid + 2) + "..." + str.substring(str.length - wid);
  };

  const columns = [
    {
      dataIndex: "transaction",
      key: "transaction",
      width: "30%",
      render: (tx: Tx) => {
        return (
          <div>
            <div style={tx.isFromSelf ? self : other}>{shorten(tx.fromA)}</div>
            <div>
              <span style={tx.isToSelf ? self : other}>
                <ArrowRightOutlined style={{ height: "1em" }} />{" "}
                {shorten(tx.toA)}
              </span>
            </div>
          </div>
        );
      }
    },
    {
      dataIndex: "transaction",
      key: "addresses",
      width: "30%",
      render: (tx: Tx) => {
        return (
          <div>
            <div style={tx.isFromSelf ? self : other}>
              {tx.from === "" ? "ABCDE" : tx.from}
            </div>
            <div style={tx.isToSelf ? self : other}>
              {tx.to === "" ? "XYZ" : tx.to}
            </div>
          </div>
        );
      }
    },
    {
      dataIndex: "transaction",
      key: "ether",
      render: (tx: Tx) => <div>{"17.000000000000000000"}</div>
    },
    {
      dataIndex: "transaction",
      key: "count",
      render: (tx: Tx) => <div>{tx.cnt}</div>
    }
  ];

  return (
    <div>
      <div>{data[0] == "" ? "WHAT" : data[0]}</div>
      <Table
        bordered={true}
        pagination={false}
        showHeader={false}
        dataSource={dataSource}
        columns={columns}
      />
    </div>
  );
};

interface Tx {
  fromA: string;
  from: string;
  toA: string;
  to: string;
  cnt: string;
  isFromSelf: boolean;
  isToSelf: boolean;
}

function str2Tx(column: string, currAddr: string): Tx {
  const parts = column.split("|");
  const tx: Tx = {
    fromA: parts[0] ?? "",
    from: parts[1] ?? "",
    toA: parts[2] ?? "",
    to: parts[3] ?? "",
    cnt: parts[4] ?? "",
    isFromSelf: (parts[0] ?? "") === currAddr,
    isToSelf: (parts[2] ?? "") === currAddr
  };
  return tx;
}

const other = {
  fontFamily: '"Courier New", monospace',
  fontWeight: "light",
  width: "100%",
  color: "#222222" /* "#f79090", */,
  overflow: "auto"
};

const self = {
  fontFamily: '"Courier New", monospace',
  fontWeight: "bold",
  width: "100%",
  color: "darkpurple",
  overflow: "auto"
};
