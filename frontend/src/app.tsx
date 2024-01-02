import React, { useState, useEffect } from "react";
import { Export, Reload } from "../wailsjs/go/main/App";
import * as Wails from "../wailsjs/runtime";
import "./App.css";
import {
  FooterDiv as Footer,
  HeaderDiv as Header,
  SideBar,
} from "./components";
import { BarChart } from "./components/barchart";

export const App: React.FC = () => {
  const [address, setAddress] = useState("trueblocks.eth");
  const [status, setStatus] = useState("Enter an address at the left...");

  const mode = "";
  const exportTxs = async () => {
    if (status == "Loading...") {
      return;
    }
    setStatus("Loading...");
    await Export(address, mode, false);
    setStatus("");
  };

  const reloadTxs = async () => {
    if (status == "Loading...") {
      return;
    }
    setStatus("Loading...");
    await Reload(address, mode, false);
    setStatus("");
  };

  return (
    <div className="panel window">
      <Header />
      <div className="panel main">
        <SideBar
          address={address}
          setAddress={setAddress}
          exportTxs={exportTxs}
          reloadTxs={reloadTxs}
        />
        <InnerPanel />
      </div>
      <Footer />
    </div>
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
  var [info, setInfo] = useState("This is the info panel.");
  var [years, setYears] = useState("This is the years panel");
  var [months, setMonths] = useState("This is the months panel");
  var [toCount, setToCount] = useState("This is the toCount panel");
  var [fromCount, setFromCount] = useState("This is the fromCount panel");
  var [fromTopTen, setFromTopTen] = useState("This is the fromTopTen panel");
  var [toTopTen, setToTopTen] = useState("This is the toTopTen panel");

  useEffect(() => {
    const update = (info: string) => {
      setInfo(info);
    };
    Wails.EventsOn("info", update);
    return () => {
      Wails.EventsOff("info");
    };
  }, []);

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
    const update = (toCount: string) => {
      setToCount(toCount);
    };
    Wails.EventsOn("toCount", update);
    return () => {
      Wails.EventsOff("toCount");
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
    const update = (toTopTen: string) => {
      setToTopTen(toTopTen);
    };
    Wails.EventsOn("toTopTen", update);
    return () => {
      Wails.EventsOff("toTopTen");
    };
  }, []);

  let stringArray: string[] = [
    info,
    years,
    months,
    fromTopTen,
    fromCount,
    toTopTen,
    toCount,
  ];

  var info1 = function () {
    return (
      <div className="panel inner-panel-body-single">
        <BorderedWord word="info" content={info} />
      </div>
    );
  };

  var dates = function () {
    return (
      <div className="panel inner-panel-body-triple">
        <BorderedWord word="Apps per Year" content={<BarChart str={years} />} />
        <BorderedWord
          word="Apps per Month"
          content={<BarChart str={months} />}
        />
      </div>
    );
  };

  var from = function () {
    const columns = fromTopTen.split(",");
    return (
      <div className="panel inner-panel-body-triple">
        <BorderedWord
          word="From"
          content={columns.map((column, index) => (
            <div>{column.trim()}</div>
          ))}
        />
        {/* <BorderedWord word="fromCount" content={fromCount} /> */}
      </div>
    );
  };

  var to = function () {
    const columns = toTopTen.replace("/,$/g", "").split(",");
    return (
      <div className="panel inner-panel-body-triple">
        <BorderedWord
          word="To"
          content={columns.map((column, index) => (
            <div>{column.trim()}</div>
          ))}
        />
        {/* <BorderedWord word="toCount" content={toCount} /> */}
      </div>
    );
  };

  return (
    <div className="inner-panel-body">
      {info1()}
      {dates()}
      {from()}
      {to()}
    </div>
  );
};

type Props = {
  str: string;
};

const CommaSeparatedDivs: React.FC<Props> = ({ str }) => {
  const columns = str.split(",");
  return (
    <div>
      {columns.map((column, index) => (
        <div>{column.trim()}</div>
      ))}
    </div>
  );
};

const Logger = () => {
  var [progress, setProgress] = useState("");
  var [error, setError] = useState("");
  useEffect(() => {
    const update = (progress: string) => {
      setProgress(progress);
    };
    Wails.EventsOn("progress", update);
    return () => {
      Wails.EventsOff("progress");
    };
  }, []);

  useEffect(() => {
    const update = (error: string) => {
      setError(error);
    };
    Wails.EventsOn("error", update);
    return () => {
      Wails.EventsOff("error");
    };
  }, []);

  var classStr = "panel inner-panel-footer";
  var content = progress;
  if (error != "") {
    classStr += " error";
    content = error;
  } else if (progress == "") {
    classStr += " empty";
    content = "Enter new addresses to the left...";
  }

  return (
    <div className={classStr}>
      <pre>{content}</pre>
    </div>
  );
};

interface BorderedWordProps {
  word: string;
  content: any;
}

const BorderedWord: React.FC<BorderedWordProps> = ({ word, content }) => {
  return (
    <div className="bordered-word-container">
      <div className="bordered-word">{word}</div>
      <div className="bordered-word-content">{content}</div>
    </div>
  );
};

// type Props1 = {
//   str: string;
// };

// const BarChart: React.FC<Props1> = ({ str }) => {
//   if (!str || !str.includes(",")) {
//     return <div>Loading...</div>;
//   }
//   str = str.replace(/^,/, "").replace(/,$/, "");

//   const parseDataString = (
//     str: string
//   ): { labels: string[]; data: number[] } => {
//     const pairs = str.split(",");
//     const labels: string[] = [];
//     const data: number[] = [];
//     pairs.forEach((pair) => {
//       const [label, count] = pair.split("-").map((str) => str.trim());
//       labels.push(label);
//       data.push(parseInt(count));
//     });

//     return { labels, data };
//   };

//   const { labels, data } = parseDataString(str);

//   const generateColors = (length: number): string[] => {
//     return Array.from(
//       { length },
//       (_, i) => `hsl(${(i / length) * 360}, 70%, 50%)`
//     );
//   };

//   const chartData: ChartData<"bar"> = {
//     labels: labels,
//     datasets: [
//       {
//         label: "",
//         data: data,
//         backgroundColor: generateColors(data.length),
//       },
//     ],
//   };

//   const options: ChartOptions<"bar"> = {
//     plugins: {
//       tooltip: {
//         callbacks: {
//           label: function (context) {
//             let label = context.dataset.label || "";
//             if (label) {
//               label += ": ";
//             }
//             if (context.parsed.y !== null) {
//               label += new Intl.NumberFormat().format(context.parsed.y);
//             }
//             return label;
//           },
//         },
//       },
//       legend: {
//         display: false,
//       },
//     },
//     scales: {
//       y: {
//         ticks: {
//           callback: function (value) {
//             return new Intl.NumberFormat().format(value as number);
//           },
//         },
//       },
//     },
//   };

//   return <Bar data={chartData} options={options} />;
// };
