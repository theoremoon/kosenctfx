import {
  Button,
  Code,
  Modal,
  ModalBody,
  ModalContent,
  ModalOverlay,
  Stack,
  Switch,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
  useDisclosure,
} from "@chakra-ui/react";
import Loading from "components/loading";
import TaskModalBody from "components/taskmodalbody";
import useMessage from "lib/useMessage";
import React, { useCallback, useState } from "react";
import Right from "../../../components/right";
import { api } from "../../../lib/api";
import useAdminTasks, { Task } from "../../../lib/api/admin/tasks";
import { useRouter } from "next/router";
import AdminLayout from "components/adminLayout";

type taskElementProps = {
  task: Task;
  isOpened: boolean;
  onUpdateOpened: (opened: boolean) => void;
  onClickCallback: () => void;
};

const TaskElement = ({
  task,
  isOpened,
  onUpdateOpened,
  onClickCallback,
  ...props
}: taskElementProps) => {
  return (
    <>
      <Td title={task.name}>{task.name}</Td>
      <Td>{task.score}</Td>
      <Td>{task.solved_by.length}</Td>
      <Td>{task.category}</Td>
      <Td>{task.author}</Td>
      <Td>
        <Switch
          isChecked={isOpened}
          onChange={() => {
            onUpdateOpened(!isOpened);
          }}
        />
      </Td>
      <Td>{task.is_survey && "🗒️"}</Td>
      <Td style={{ cursor: "pointer" }} onClick={onClickCallback}>
        👀
      </Td>
      <Td>
        <Code colorScheme="whiteAlpha" title={task.flag}>
          {task.flag}
        </Code>
      </Td>
    </>
  );
};

interface TasksProps {
  tasks: Task[];
}

const Tasks = ({ tasks }: TasksProps) => {
  const { mutate } = useAdminTasks();
  const { message, error } = useMessage();

  const openState = new Map(tasks.map((t) => [t.id, t.is_open]));
  const [taskOpenState, setTaskOpenState] = useState(openState);
  const router = useRouter();

  const openCloseChallenge = useCallback(() => {
    tasks.forEach(async (t) => {
      const next_state = taskOpenState.get(t.id);
      if (t.is_open === next_state) {
        return;
      }

      const endpoint = next_state
        ? "/admin/open-challenge"
        : "/admin/close-challenge";
      try {
        const res = await api.post(endpoint, {
          name: t.name,
        });
        message(res);
      } catch (e) {
        error(e);
      }
    });
    mutate();
  }, [api, message, error, tasks]);

  return (
    <AdminLayout>
        <Table maxW="100%" size="sm">
          <Thead>
            <Tr>
              <Th>Name</Th>
              <Th>Score</Th>
              <Th>#Solve</Th>
              <Th>Category</Th>
              <Th>Author</Th>
              <Th>Is Open?</Th>
              <Th>Is Survey?</Th>
              <Th>Preview</Th>
              <Th>Flag</Th>
            </Tr>
          </Thead>
          <Tbody>
            {tasks.map((task) => (
              <Tr key={task.name}>
                <TaskElement
                  task={task}
                  isOpened={taskOpenState.get(task.id) || false}
                  onUpdateOpened={(isOpen) => {
                    setTaskOpenState((prev) =>
                      new Map(prev).set(task.id, isOpen)
                    );
                  }}
                  onClickCallback={() => {
                    router.push(`/admin/tasks/${task.id}`, undefined, {
                      shallow: true,
                      scroll: false,
                    });
                  }}
                />
              </Tr>
            ))}
          </Tbody>
        </Table>

        <Right>
          <Button onClick={openCloseChallenge}>Open/Close Challenges</Button>
        </Right>
    </AdminLayout>
  );
};

const TasksDefault = () => {
  const { data: tasks } = useAdminTasks();

  if (!tasks) {
    return <Loading />;
  }
  return <Tasks tasks={tasks} />;
};

export default TasksDefault;
