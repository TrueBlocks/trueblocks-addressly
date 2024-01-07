import React from "react";
import { Bar } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ChartData,
  ChartOptions
} from "chart.js";

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);

export interface Props {
  dataStr: string;
  clickHandler: (chartType: string) => void;
}

export const BarChart: React.FC<Props> = ({ dataStr, clickHandler }) => {
  if (dataStr.length === 0) {
    return <div>Loading...</div>;
  }
  dataStr = dataStr.replace(/^,/, "").replace(/,$/, "");

  const { labels, data } = parseDataString(dataStr);

  const generateColors = (length: number): string[] => {
    return Array.from(
      { length },
      (_, i) => `hsl(${(i / length) * 360}, 70%, 50%)`
    );
  };

  const chartData: ChartData<"bar"> = {
    labels,
    datasets: [
      {
        label: "",
        data,
        backgroundColor: generateColors(data.length)
      }
    ]
  };

  const options: ChartOptions<"bar"> = {
    plugins: {
      tooltip: {
        callbacks: {
          label: function (context) {
            let label = context.dataset.label || "";
            if (label.length > 0) {
              label += ": ";
            }
            if (context.parsed.y !== null) {
              label += new Intl.NumberFormat().format(context.parsed.y);
            }
            return label;
          }
        }
      },
      legend: {
        display: false
      }
    },
    scales: {
      y: {
        ticks: {
          callback: function (value) {
            return new Intl.NumberFormat().format(value as number);
          }
        }
      }
    }
  };

  const handleBarClick = (elements: any[]) => {
    if (elements.length > 0) {
      const firstElement = elements[0];
      // console.log(
      //   `Bar clicked: ${firstElement.index} ${labels[firstElement.index]} ${
      //     data[firstElement.index]
      //   }`
      // );
      if (labels[firstElement.index].length == 4) {
        clickHandler("month");
      } else {
        clickHandler("year");
      }
    }
  };

  return (
    <Bar
      data={chartData}
      options={options}
      onClick={(event: React.MouseEvent<HTMLCanvasElement>) => {
        const canvas = event.currentTarget;
        const chart = ChartJS.getChart(canvas);
        if (chart) {
          const elements = chart.getElementsAtEventForMode(
            event.nativeEvent,
            "nearest",
            { intersect: true },
            false
          );
          handleBarClick(elements);
        }
      }}
    />
  );
};

const parseDataString = (str: string): { labels: string[]; data: number[] } => {
  const pairs = str.split(",");
  const labels: string[] = [];
  const data: number[] = [];
  pairs.forEach((pair) => {
    const [label, count] = pair.split("|").map((str) => str.trim());
    labels.push(label);
    data.push(parseInt(count));
  });

  return { labels, data };
};
