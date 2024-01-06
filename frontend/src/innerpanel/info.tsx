import React, { useContext } from "react";
import "./inner.css";
import { Card } from "antd";
import { Table, Typography } from "antd";
import { AppContext } from "../appcontext";
import { EditableAddress } from "../components";
const { Text } = Typography;

export const Info: React.FC = () => {
  var { info } = useContext(AppContext);
  return (
    <Card className="inner-panel-body-single" style={{ width: "100%" }}>
      {/* <div style={{ color: "black" }}>{info}</div> */}
      <MyTableComponent />
    </Card>
  );
};

interface TableRow {
  key: number;
  field1: React.ReactElement<any, any>;
  field2?: string;
  field3?: React.ReactElement<any, any>;
  field4?: string;
  field5?: React.ReactElement<any, any>;
  field6?: string;
}

const MyTableComponent: React.FC = () => {
  var { info } = useContext(AppContext);

  const fields = info.split(",");
  const fieldNames = [
    "Name",
    "Address",
    "Appearance Count",
    "First Appearance",
    "Latest Appearance",
    "ETH Balance",
    "Block Span",
    "Blocks Between",
    "USD Balance",
  ];

  const combinedArray = fieldNames.reduce(
    (acc, name, index) => [...acc, name, fields[index] ?? ""],
    [] as string[]
  );

  const tableData: TableRow[] = [];
  for (let i = 0; i < combinedArray.length; i += 6) {
    const row: TableRow = {
      key: i,
      field1: (
        <Text strong style={{ fontSize: "1.1em", textAlign: "right" }}>
          {combinedArray[i + 0] + ":" ?? ""}
        </Text>
      ),
      field2: combinedArray[i + 1] + "|" + combinedArray[i + 3] ?? "",

      field3: (
        <Text strong style={{ fontSize: "1.1em", textAlign: "right" }}>
          {combinedArray[i + 2] + ":" ?? ""}
        </Text>
      ),
      field4: combinedArray[i + 3] ?? "",

      field5: (
        <Text strong style={{ fontSize: "1.1em", textAlign: "right" }}>
          {combinedArray[i + 4] + ":" ?? ""}
        </Text>
      ),
      field6: combinedArray[i + 5] ?? "",
    };
    tableData.push(row);
  }

  const widths = ["10%", "10%", "10%", "20%", "10%", "20%"];
  const columns = [
    {
      dataIndex: "field1",
      key: "field1",
      width: widths[0],
    },
    {
      dataIndex: "field2",
      key: "field2",
      width: widths[1],
      maxWidth: widths[1],
      render: (text: string, record: TableRow, index: number) => {
        if (index === 0) {
          return (
            <span style={{ fontSize: "1.2em", color: "dodgerblue" }}>
              <EditableAddress name={text} />
            </span>
          );
        }
        return text;
      },
      ellipsis: true,
    },
    {
      dataIndex: "field3",
      key: "field3",
      width: widths[2],
    },
    {
      dataIndex: "field4",
      key: "field4",
      width: widths[3],
      ellipsis: true,
    },
    {
      dataIndex: "field5",
      key: "field5",
      width: widths[4],
    },
    {
      dataIndex: "field6",
      key: "field6",
      width: widths[5],
      render: (text: string, record: TableRow, index: number) => {
        if (index === 1 || index === 2) {
          return (
            <span style={{ fontSize: "1.2em", color: "dodgerblue" }}>
              {text}
            </span>
          );
        }
        return text;
      },
      ellipsis: true,
    },
  ];

  return (
    <Table
      dataSource={tableData}
      columns={columns}
      showHeader={false}
      pagination={false}
      size="small"
      style={
        {
          whiteSpace: "nowrap",
          "--table-padding-vertical": "4px",
          "--table-padding-horizontal": "8px",
          border: "1px solid lightgray",
        } as React.CSSProperties
      }
      rowKey={(record) => record.key.toString()}
    />
  );
};
