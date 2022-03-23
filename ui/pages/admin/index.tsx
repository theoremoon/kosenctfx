import { ErrorMessage } from "@hookform/error-message";
import InlineFormControl from "components/inlineformcontrol";
import Loading from "components/loading";
import Right from "components/right";
import ReactECharts from "echarts-for-react";
import useMessage from "lib/useMessage";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import Area from "../../components/area";
import { api } from "../../lib/api";
import useConfig, { Config } from "../../lib/api/admin/config";
import useScoreboard from "../../lib/api/scoreboard";
import { dateFormat, parseDateString } from "../../lib/date";
import FormWrapper from "../../components/formwrapper";
import FormItem from "../../components/formitem";
import Label from "../../components/label";
import Input from "../../components/input";
import InlineLabel from "../../components/inlinelabel";
import InlineInput from "../../components/inlineinput";
import Button from "../../components/button";
import Divider from "../../components/divider";

type ConfigParams = {
  start_at: string;
  end_at: string;
  ctf_open: boolean;
  register_open: boolean;
  lock_second: number;
  lock_count: number;
  lock_duration: number;
  score_expr: string;
};

type ScoreEmulateResponse = number[];

interface AdminConfigImplProps {
  data: Config;
}

const AdminConfigImpl = ({ data }: AdminConfigImplProps) => {
  const { mutate } = useConfig();
  const { start_at: _start_at, end_at: _end_at, ..._data } = data;
  const { data: scoreboard } = useScoreboard();

  const params: ConfigParams = {
    start_at: dateFormat(data.start_at),
    end_at: dateFormat(data.end_at),
    ..._data,
  };

  const {
    register,
    handleSubmit,
    getValues,
    formState: { errors },
  } = useForm<ConfigParams>({ defaultValues: params });

  const { message, error } = useMessage();

  const onSubmit = async (data: ConfigParams) => {
    const { start_at, end_at, ..._data } = data;
    const newConfig = {
      start_at: parseDateString(data.start_at),
      end_at: parseDateString(data.end_at),
      ..._data,
    };

    try {
      const res = await api.post("/admin/set-config", newConfig);
      message(res);
      mutate();
    } catch (e) {
      error(e);
    }
  };

  // スコア遷移グラフまわり
  const [scoreGraph, setScoreGraph] = useState<number[]>([]);
  const [numberOfSolves, setNumberOfSolves] = useState(200);
  const scoreEmulate = () => {
    api
      .get<ScoreEmulateResponse>("/admin/score-emulate", {
        params: {
          maxCount: numberOfSolves,
          expr: getValues("score_expr"),
        },
      })
      .then((r) => {
        setScoreGraph(r.data);
      });
  };

  const getScoreboard = () => {
    const link = document.createElement("a");
    link.href =
      "data:applicaion/json;charset=utf-8," +
      encodeURIComponent(JSON.stringify(scoreboard));
    link.download = "scoreboard.json";
    document.body.appendChild(link);
    link.click();
    setTimeout(() => {
      if (link.parentNode) {
        link.parentNode.removeChild(link);
      }
    }, 1000);
  };

  // score seriesを全部再計算する恐怖のメソッド
  const recalc = async () => {
    try {
      const res = await api.post("/admin/recalc-series");
      message(res);
    } catch (e) {
      error(e);
    }
  };

  return (
    <div className="w-min md:w-2/3 mx-auto mt-10">
      <form onSubmit={handleSubmit(onSubmit)}>
        <ErrorMessage errors={errors} name="Error" />

        <FormItem>
          <InlineLabel htmlFor="start_at">CTF Starts at</InlineLabel>
          <InlineInput
            id="start_at"
            type="text"
            {...register("start_at", {
              required: true,
              validate: (value: string) => parseDateString(value) !== null,
            })}
          />
        </FormItem>

        <FormItem>
          <InlineLabel htmlFor="end_at">CTF Ends at</InlineLabel>
          <InlineInput
            id="end_at"
            type="text"
            {...register("end_at", {
              required: true,
              validate: (value: string) => parseDateString(value) !== null,
            })}
          />
        </FormItem>

        <FormItem>
          <InlineLabel htmlFor="ctf_open">CTF is open</InlineLabel>
          <input
            {...register("ctf_open")}
            type="checkbox"
            className="w-6 h-6 text-pink-600"
          />

          <InlineLabel htmlFor="register_open">
            Registration is open
          </InlineLabel>
          <input
            {...register("register_open")}
            type="checkbox"
            className="w-6 h-6 text-pink-600"
          />
        </FormItem>

        <FormItem>
          <InlineLabel>Submission lock second</InlineLabel>
          <InlineInput type="number" {...register("lock_second")} />

          <InlineLabel>Lock trigger count</InlineLabel>
          <InlineInput type="number" {...register("lock_count")} />

          <InlineLabel>Lock trigger duration</InlineLabel>
          <InlineInput type="number" {...register("lock_duration")} />
        </FormItem>

        <Area>
          <FormItem>
            <InlineLabel htmlFor="end_at">Max Solves</InlineLabel>
            <InlineInput
              value={numberOfSolves}
              onChange={(e) => setNumberOfSolves(Number(e.target.value))}
            />
          </FormItem>
          <Right>
            <Button onClick={scoreEmulate}>Draw Graph</Button>
          </Right>
          <ReactECharts
            style={{ flex: 1 }}
            option={{
              tooltip: {
                trigger: "axis",
              },
              xAxis: {
                name: "number of solves",
                nameLocation: "center",
                type: "category",
                data: scoreGraph.map((_, i) => i),
              },
              yAxis: {
                name: "score",
                type: "value",
              },
              series: [
                {
                  data: scoreGraph,
                  type: "line",
                },
              ],
            }}
          />
        </Area>

        <FormItem>
          <Right>
            <Button type="submit">Apply Changes</Button>
          </Right>
        </FormItem>
      </form>

      <Divider />

      <div className="grid grid-cols-3 gap-x-4">
        <Button
          onClick={() => {
            getScoreboard();
          }}
        >
          Scoreboard for CTFtime
        </Button>

        <Button onClick={() => {}}>Scoreboard for CTF-ratings</Button>

        <Button
          onClick={() => {
            recalc();
          }}
        >
          Recalc All Score Series
        </Button>
      </div>
    </div>
  );
};

const AdminConfig = () => {
  const { data } = useConfig();

  if (!data) {
    return <Loading />;
  }
  return <AdminConfigImpl data={data} />;
};

export default AdminConfig;
