import AdminLayout from "components/adminLayout";
import useAdminTasks, { Task } from "lib/api/admin/tasks";
import Loading from "components/loading";
import useAgents, { Agent } from "lib/api/admin/agent";
import { useForm, SubmitHandler } from "react-hook-form";
import useMessage from "lib/useMessage";
import { api } from "lib/api";
import useLivingDeployments, { Deployment } from "lib/api/admin/deployment";
import { KeyedMutator } from "swr";

interface CreateNewDeploymentProps {
  tasks: Task[];
  agents: Agent[];
  mutateDeployments: KeyedMutator<Deployment[]>;
}

type NewDeploymentParams = {
  task_id: number;
  agent_id: string;
};

const CreateNewDeployment = ({
  tasks,
  agents,
  mutateDeployments,
}: CreateNewDeploymentProps) => {
  const { register, handleSubmit } = useForm<NewDeploymentParams>();
  const { message, error } = useMessage();
  const onSubmit: SubmitHandler<NewDeploymentParams> = async (values) => {
    try {
      console.log(values);
      const res = await api.post("/admin/request-deploy", {
        task_id: +values.task_id,
        agent_id: values.agent_id,
      });
      message(res);
      mutateDeployments();
    } catch (e) {
      error(e);
    }
  };

  return (
    <>
      <form onSubmit={handleSubmit(onSubmit)}>
        <div className="form-group row">
          <label className="col-sm-2 col-form-label">Task</label>
          <div className="col-sm-10">
            <select
              className="form-select"
              {...register("task_id", { required: true })}
            >
              {tasks
                .filter((t) => t.deployment !== "")
                .map((t) => (
                  <option key={t.id} value={t.id}>
                    {t.name} ({t.deployment})
                  </option>
                ))}
            </select>
          </div>
        </div>

        <div className="form-group row">
          <label className="col-sm-2 col-form-label">Agent</label>
          <div className="col-sm-10">
            <select
              className="form-select"
              {...register("agent_id", { required: true })}
            >
              {agents.map((a) => (
                <option key={a.agent_id} value={a.agent_id}>
                  {a.agent_id} ({a.public_ip})
                </option>
              ))}
            </select>
          </div>
        </div>

        <button type="submit" className="btn btn-primary">
          Deploy
        </button>
      </form>
    </>
  );
};

interface DeploymentsProps {
  tasks: Task[];
  agents: Agent[];
  deployments: Deployment[];
  mutateDeployments: KeyedMutator<Deployment[]>;
}

const Deployments = ({
  tasks,
  agents,
  deployments,
  mutateDeployments,
}: DeploymentsProps) => {
  return (
    <>
      <h5 className="mt-4">Current Running Deployments</h5>
      <table className="table table-responsive">
        <thead>
          <tr>
            <th>ID</th>
            <th>Challenge</th>
            <th>Agent</th>
            <th>Status</th>
            <th>Port</th>
            <th>Requested At</th>
            <th>Retires At</th>
          </tr>
        </thead>
        <tbody>
          {deployments.map((d) => (
            <tr key={d.id}>
              <td>{d.id}</td>
              <td>{d.challenge?.name || "#"}</td>
              <td>{d.agent?.agent_id || "#"}</td>
              <td>{d.status}</td>
              <td>{d.port}</td>
              <td>{d.requested_at}</td>
              <td>{d.retires_at}</td>
            </tr>
          ))}
        </tbody>
      </table>

      <h5 className="mt-4">Create New Deployment</h5>
      <CreateNewDeployment
        tasks={tasks}
        agents={agents}
        mutateDeployments={mutateDeployments}
      />
    </>
  );
};

const DeploymentsDefault = () => {
  const { data: tasks } = useAdminTasks();
  const { data: agents } = useAgents();
  const { data: deployments, mutate: mutateDeployments } =
    useLivingDeployments();

  if (
    tasks === undefined ||
    agents === undefined ||
    deployments === undefined
  ) {
    return <Loading />;
  }
  return (
    <Deployments
      tasks={tasks}
      agents={agents}
      deployments={deployments}
      mutateDeployments={mutateDeployments}
    />
  );
};

DeploymentsDefault.getLayout = AdminLayout;

export default DeploymentsDefault;
