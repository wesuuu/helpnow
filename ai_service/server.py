import grpc
from concurrent import futures
import time
import sys
import os

# Add gen directory to path to import generated files
sys.path.append(os.path.join(os.path.dirname(__file__), 'gen'))

# Try importing generated code. If not found, these will fail, 
# preventing the server from starting until protoc is run.
try:
    import ai_service_pb2
    import ai_service_pb2_grpc
except ImportError:
    print("Warning: Generated Protobuf files not found. Run 'make proto' first.")
    # Create dummy classes to allow syntax checking
    class ai_service_pb2_grpc:
        AIServiceServicer = object
        def add_AIServiceServicer_to_server(x, y): pass
    class ai_service_pb2:
        ExecuteRoutineResponse = object

from agents.marketing import MarketingAgentWrapper

class AIService(ai_service_pb2_grpc.AIServiceServicer):
    def __init__(self):
        pass

    def ExecuteRoutine(self, request, context):
        print(f"Received ExecuteRoutine request: {request.routine_id}")
        
        # Parse inputs
        inputs = dict(request.input_params)
        model_config = inputs.get('model_config', '{}')
        # TODO: Parse model_config JSON string
        
        # Initialize Agent with specific config
        # Ideally we cache agents or manage them, but for now we instantiate per request
        agent_wrapper = MarketingAgentWrapper(model_config=model_config)
        
        workflow = inputs.get('workflow', f"Execute routine {request.routine_id}")

        try:
            result = agent_wrapper.run_routine(workflow, inputs)
            status = "COMPLETED"
        except Exception as e:
            print(f"Error executing routine: {e}")
            result = str(e)
            status = "FAILED"

        # Generate a mock execution ID
        execution_id = f"exec-{int(time.time())}"

        return ai_service_pb2.ExecuteRoutineResponse(
            execution_id=execution_id,
            status=status
        )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    ai_service_pb2_grpc.add_AIServiceServicer_to_server(AIService(), server)
    server.add_insecure_port('[::]:50051')
    print("AI Service started on port 50051")
    server.start()
    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == '__main__':
    serve()
