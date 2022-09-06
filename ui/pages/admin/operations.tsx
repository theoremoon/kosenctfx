import AdminLayout from "components/adminLayout";
import { api } from "lib/api";
import useConfig from "lib/api/admin/config";
import useScoreboard from "lib/api/scoreboard";
import useTasks from "lib/api/tasks";
import useMessage from "lib/useMessage";

const AdminOperations = () => {
  const { data: scoreboard } = useScoreboard([]);
  const { data: tasks } = useTasks([]);
  const { message, error } = useMessage();

  // score seriesを全部再計算する恐怖のメソッド
  const recalc = async () => {
    try {
      const res = await api.post("/admin/recalc-series");
      message(res);
    } catch (e) {
      error(e);
    }
  };

  const getScoreboard = () => {
    const link = document.createElement("a");
    link.href =
      "data:applicaion/json;charset=utf-8," +
      encodeURIComponent(
        JSON.stringify({
          tasks: tasks?.map((task) => task.name),
          standings: scoreboard?.map((team) => ({
            pos: team.pos,
            team: team.team,
            score: team.score,
            taskStats: team.taskStats,
            lastAccept: team.last_submission,
          })),
        })
      );
    link.download = "scoreboard.json";
    document.body.appendChild(link);
    link.click();
    setTimeout(() => {
      if (link.parentNode) {
        link.parentNode.removeChild(link);
      }
    }, 1000);
  };

  const { data: config, error: configError } = useConfig();
  if (config === undefined || configError !== undefined) {
    return <></>;
  }

  return (
    <>
      <h5 className="mt-4">CTFtime Scoreboard</h5>
      <p>Get a CTFtime-compatibility scoreboard json file</p>
      <button
        type="button"
        className="btn btn-primary mx-2"
        onClick={() => getScoreboard()}
      >
        Download CTFTime scoreboard
      </button>

      <h5 className="mt-4">Recalc series</h5>
      <p>
        Trigger to recalculate all team score series. You need to use this
        feature after you changed the score expr
      </p>
      <button
        type="button"
        className="btn btn-primary mx-2"
        onClick={() => recalc()}
      >
        Recalc all series
      </button>
    </>
  );
};

AdminOperations.getLayout = AdminLayout;

export default AdminOperations;
