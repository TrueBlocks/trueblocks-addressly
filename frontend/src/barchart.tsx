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
  ChartOptions,
} from "chart.js";

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);

interface BarChartProps {
  dataString: string;
}

export const BarChart: React.FC<BarChartProps> = ({ dataString }) => {
  if (!dataString || !dataString.includes(",")) {
    return <div>Invalid dataString: {dataString}</div>;
  }
  dataString = dataString.replace(/^,/, "").replace(/,$/, "");

  const parseDataString = (
    dataStr: string
  ): { labels: string[]; data: number[] } => {
    const pairs = dataStr.split(",");
    const labels: string[] = [];
    const data: number[] = [];
    pairs.forEach((pair) => {
      const [label, count] = pair.split("-").map((str) => str.trim());
      labels.push(label);
      data.push(parseInt(count));
    });

    return { labels, data };
  };

  const { labels, data } = parseDataString(dataString);

  const generateColors = (length: number): string[] => {
    return Array.from(
      { length },
      (_, i) => `hsl(${(i / length) * 360}, 70%, 50%)`
    );
  };

  const chartData: ChartData<"bar"> = {
    labels: labels,
    datasets: [
      {
        label: "",
        data: data,
        backgroundColor: generateColors(data.length),
      },
    ],
  };

  const options: ChartOptions<"bar"> = {
    plugins: {
      tooltip: {
        callbacks: {
          label: function (context) {
            let label = context.dataset.label || "";
            if (label) {
              label += ": ";
            }
            if (context.parsed.y !== null) {
              label += new Intl.NumberFormat().format(context.parsed.y);
            }
            return label;
          },
        },
      },
      legend: {
        display: false,
      },
    },
    scales: {
      y: {
        ticks: {
          callback: function (value) {
            return new Intl.NumberFormat().format(value as number);
          },
        },
      },
    },
  };

  return <Bar data={chartData} options={options} />;
};
