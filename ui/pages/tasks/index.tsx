import Loading from "../../components/loading";
import Tasks from "../../components/tasks";
import useAccount from "../../lib/api/account";
import useTasks from "../../lib/api/tasks";

const TasksDefault = () => {
  const { data: tasks } = useTasks();
  const { data: account } = useAccount();
  if (!tasks || account === undefined) {
    return <Loading />;
  }
  return <Tasks tasks={tasks} team={account} />;
};

export default TasksDefault;
