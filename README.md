“Async4Ai” is a Go framework of asynchronously scheduling tasks for accelerating data process. It's used under the case of continuously working on multi-stage tasks.

- Why not use Message Queue(MQ) directly:
   - Message Queue can only distribute tasks
   - Async4Ai can also update/inqury tasks via offering web API within web component which allow visits from outside and into underlying DB.

- Why not use other task manager tool like celery and machinery directly:
   - Celery and machinery are comprehensive and powerful tools not only dealing with async task.
   - Async4Ai are more light-weighted and specifically designed for deep learning tasks in optimized for group workflow.

