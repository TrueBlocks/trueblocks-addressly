import React, { useEffect, useContext } from "react";
import * as Wails from "../../wailsjs/runtime";
import { Table, Typography } from "antd";
import { AppContext } from "../appcontext";
const { Text } = Typography;

interface TableRow {
  key: number;
  field1: React.ReactElement<any, any>;
  field2?: string;
  field3?: React.ReactElement<any, any>;
  field4?: string;
  field5?: React.ReactElement<any, any>;
  field6?: string;
}

export const MyTableComponent: React.FC = () => {
  var { info, setInfo } = useContext(AppContext);

  const fields = info.split(",");
  const fieldNames = [
    "Your query",
    "Address",
    "Appearance Count",
    "First Appearance",
    "First Timestamp",
    "First Date",
    "Latest Appearance",
    "Latest Timestamp",
    "Latest Date",
    "Lifetime Span",
    "Blocks Between",
    "ETH Balance",
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
      field2: combinedArray[i + 1] ?? "",

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
      overflow: "hidden",
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
      overflow: "hidden",
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
      overflow: "hidden",
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
