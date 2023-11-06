import AdminLayout from "components/adminLayout";
import { ErrorMessage } from "@hookform/error-message";
import ReactECharts from "echarts-for-react";
import useMessage from "lib/useMessage";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { api } from "lib/api";
import useConfig, { Config } from "lib/api/admin/config";
import { dateFormat, parseDateString } from "lib/date";

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

type AdminConfigProps = {
  config: Config;
};

const AdminConfig = ({ config: defaultConfig }: AdminConfigProps) => {
  const { data: config, mutate } = useConfig();
  const params: ConfigParams = {
    ...(config || defaultConfig),
    start_at: dateFormat(config?.start_at || 0),
    end_at: dateFormat(config?.end_at || 0),
  };

  const {
    register,
    handleSubmit,
    getValues,
    formState: { errors },
  } = useForm<ConfigParams>({ defaultValues: params });

  const { message, error } = useMessage();

  const onSubmit = async (data: ConfigParams) => {
    const newConfig = {
      ...data,
      start_at: parseDateString(data.start_at),
      end_at: parseDateString(data.end_at),
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

  return (
    <>
      <form onSubmit={handleSubmit(onSubmit)}>
        <ErrorMessage errors={errors} name="ctf_open" />

        <h5 className="mt-4">CTF Config</h5>
        <div className="form-group row">
          <label className="col-sm-2 col-form-label">CTF Starts at</label>
          <div className="col-sm-10">
            <input
              className="form-control"
              id="start_at"
              {...register("start_at", {
                required: true,
                validate: (value: string) => parseDateString(value) !== null,
              })}
            />
          </div>
        </div>

        <div className="form-group row">
          <label className="col-sm-2 col-form-label">CTF Ends at</label>
          <div className="col-sm-10">
            <input
              className="form-control"
              id="end_at"
              {...register("end_at", {
                required: true,
                validate: (value: string) => parseDateString(value) !== null,
              })}
            />
          </div>
        </div>

        <div className="form-group row">
          <label className="col-sm-2 col-form-label">CTF is open</label>
          <div className="col-sm-10 form-check form-switch">
            <input
              className="form-check-input"
              type="checkbox"
              role="switch"
              id="ctf_open"
              {...register("ctf_open")}
            />
          </div>
        </div>

        <div className="form-group row">
          <label className="col-sm-2 col-form-label">
            Registration is open
          </label>
          <div className="col-sm-10 form-check form-switch">
            <input
              className="form-check-input"
              type="checkbox"
              role="switch"
              id="register_open"
              {...register("register_open")}
            />
          </div>
        </div>

        <h5 className="mt-4">Submission Lock</h5>
        <div className="form-group row">
          <label className="col-sm-2 col-form-label">
            lock duration (seconds)
          </label>
          <div className="col-sm-10">
            <input
              className="form-control"
              id="lock_second"
              {...register("lock_second", { required: true })}
            />
          </div>
        </div>

        <div className="form-group row">
          <label className="col-sm-2 col-form-label">
            lock threshold count
          </label>
          <div className="col-sm-10">
            <input
              className="form-control"
              id="lock_count"
              {...register("lock_count", { required: true })}
            />
          </div>
        </div>

        <div className="form-group row">
          <label className="col-sm-2 col-form-label">
            threshold count duration
          </label>
          <div className="col-sm-10">
            <input
              className="form-control"
              id="lock_duration"
              {...register("lock_duration", { required: true })}
            />
          </div>
        </div>

        <h5 className="mt-4">Score expr</h5>
        <div className="form-floating">
          <textarea
            className="form-control"
            {...register("score_expr")}
            style={{ height: "8em" }}
          />
        </div>

        <h5 className="mt-4">Score emulate</h5>
        <div className="form-group row">
          <label className="col-sm-2 col-form-label">Max Solves</label>
          <div className="col-sm-10">
            <input
              className="form-control"
              value={numberOfSolves}
              onChange={(e) => setNumberOfSolves(Number(e.target.value))}
            />
          </div>
        </div>
        <button
          type="button"
          className="btn btn-primary"
          onClick={scoreEmulate}
        >
          Draw Graph
        </button>
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

        <button type="submit" className="btn btn-primary">
          Apply Changes
        </button>
      </form>
    </>
  );
};

const AdminConfigDefault = () => {
  const { data: config, error } = useConfig();
  if (config === undefined || error !== undefined) {
    return <></>;
  }
  return <AdminConfig config={config} />;
};

AdminConfigDefault.getLayout = AdminLayout;

export default AdminConfigDefault;
