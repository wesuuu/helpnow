import pytest
from unittest.mock import MagicMock
from agents.marketing import MarketingAgentWrapper

def test_run_routine():
    # Mock the internal Strands Agent
    agent_wrapper = MarketingAgentWrapper()
    agent_wrapper.agent = MagicMock()
    agent_wrapper.agent.run.return_value = MagicMock(content="Routine Executed Successfully")

    workflow = "Send email to client"
    inputs = {"recipient": "test@test.com"}

    result = agent_wrapper.run_routine(workflow, inputs)
    
    assert result == "Routine Executed Successfully"
    agent_wrapper.agent.run.assert_called_once()
