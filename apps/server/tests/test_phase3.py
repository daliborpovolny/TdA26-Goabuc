import re
import requests

BASE_URL = "http://localhost:80/api"  # <-- change if needed


# ---------- minimal API helpers ----------

def api_get(path):
    r = requests.get(BASE_URL + path)
    try:
        return r.json(), r
    except ValueError:
        return None, r


def api_post(path, body):
    r = requests.post(BASE_URL + path, json=body)
    try:
        return r.json(), r
    except ValueError:
        return None, r


def api_put(path, body):
    r = requests.put(BASE_URL + path, json=body)
    try:
        return r.json(), r
    except ValueError:
        return None, r


def api_delete(path):
    return requests.delete(BASE_URL + path)


# ---------- shared test state ----------

test_course_id = None
test_quiz_id = None
single_choice_question_uuid = None
multiple_choice_question_uuid = None


# ---------- tests ----------

def test_create_course_for_quizzes():
    global test_course_id

    data, response = api_post(
        "/courses",
        {
            "name": "Test Course for Phase 3",
            "description": "Testing course quizzes",
        },
    )

    assert response.ok
    assert data is not None
    assert "uuid" in data

    test_course_id = data["uuid"]


def test_list_quizzes_initially_empty():
    assert test_course_id is not None

    data, response = api_get(f"/courses/{test_course_id}/quizzes")

    assert response.ok
    assert isinstance(data, list)
    assert len(data) == 0


def test_create_quiz_single_choice():
    global test_quiz_id, single_choice_question_uuid

    quiz = {
        "title": "Introduction Quiz",
        "questions": [
            {
                "type": "singleChoice",
                "question": "What is 2 + 2?",
                "options": ["3", "4", "5", "6"],
                "correctIndex": 1,
            }
        ],
    }

    data, response = api_post(
        f"/courses/{test_course_id}/quizzes",
        quiz,
    )

    assert response.ok
    assert "uuid" in data
    assert data["title"] == quiz["title"]

    question = data["questions"][0]
    assert question["type"] == "singleChoice"

    test_quiz_id = data["uuid"]
    single_choice_question_uuid = question["uuid"]


def test_get_quiz_details():
    data, response = api_get(
        f"/courses/{test_course_id}/quizzes/{test_quiz_id}"
    )

    assert response.ok
    assert data["uuid"] == test_quiz_id
    assert len(data["questions"]) == 1


def test_update_quiz_title():
    global single_choice_question_uuid

    updated = {
        "title": "Updated Introduction Quiz",
        "questions": [
            {
                "type": "singleChoice",
                "question": "What is 2 + 2?",
                "options": ["3", "4", "5", "6"],
                "correctIndex": 1,
            }
        ],
    }

    data, response = api_put(
        f"/courses/{test_course_id}/quizzes/{test_quiz_id}",
        updated,
    )

    assert response.ok
    assert data["title"] == "Updated Introduction Quiz"

    single_choice_question_uuid = data["questions"][0]["uuid"]


def test_add_multiple_choice_question():
    global multiple_choice_question_uuid

    updated = {
        "title": "Updated Introduction Quiz",
        "questions": [
            {
                "type": "singleChoice",
                "question": "What is 2 + 2?",
                "options": ["3", "4", "5", "6"],
                "correctIndex": 1,
            },
            {
                "type": "multipleChoice",
                "question": "Which are prime numbers?",
                "options": ["2", "3", "4", "5"],
                "correctIndices": [0, 1, 3],
            },
        ],
    }

    data, response = api_put(
        f"/courses/{test_course_id}/quizzes/{test_quiz_id}",
        updated,
    )

    assert response.ok
    assert len(data["questions"]) == 2

    for q in data["questions"]:
        if q["type"] == "multipleChoice":
            multiple_choice_question_uuid = q["uuid"]


def test_submit_quiz_correct_answers():
    submission = {
        "answers": [
            {
                "uuid": single_choice_question_uuid,
                "selectedIndex": 1,
            },
            {
                "uuid": multiple_choice_question_uuid,
                "selectedIndices": [0, 1, 3],
            },
        ]
    }

    data, response = api_post(
        f"/courses/{test_course_id}/quizzes/{test_quiz_id}/submit",
        submission,
    )

    assert response.ok
    assert data["score"] == data["maxScore"]


def test_submit_quiz_incorrect_answers():
    submission = {
        "answers": [
            {
                "uuid": single_choice_question_uuid,
                "selectedIndex": 0,
            },
            {
                "uuid": multiple_choice_question_uuid,
                "selectedIndices": [0, 2],
            },
        ]
    }

    data, response = api_post(
        f"/courses/{test_course_id}/quizzes/{test_quiz_id}/submit",
        submission,
    )

    assert response.ok
    assert data["score"] < data["maxScore"]


def test_course_detail_includes_quizzes():
    data, response = api_get(f"/courses/{test_course_id}")

    assert response.ok
    assert "quizzes" in data
    assert len(data["quizzes"]) >= 1


def test_quiz_uuid_format():
    data, response = api_post(
        f"/courses/{test_course_id}/quizzes",
        {
            "title": "UUID Test Quiz",
            "questions": [
                {
                    "type": "singleChoice",
                    "question": "UUID test question?",
                    "options": ["Yes", "No"],
                    "correctIndex": 0,
                }
            ],
        },
    )

    assert response.ok

    uuid_v4 = re.compile(
        r"^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$",
        re.I,
    )

    assert uuid_v4.match(data["uuid"])


def test_delete_quiz():
    response = api_delete(
        f"/courses/{test_course_id}/quizzes/{test_quiz_id}"
    )

    assert response.status_code == 204

    _, get_response = api_get(
        f"/courses/{test_course_id}/quizzes/{test_quiz_id}"
    )

    assert get_response.status_code == 404


def test_cleanup_course():
    response = api_delete(f"/courses/{test_course_id}")
    assert response.status_code == 204
