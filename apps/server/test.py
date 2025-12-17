import requests
import uuid
import json

BASE_URL = "http://localhost/api"

session = requests.Session()
session.headers.update({
    "Content-Type": "application/json",
    "Accept": "application/json",
})

def log(title):
    print(f"\n=== {title} ===")

def check(resp):
    print(resp.status_code)
    if resp.content:
        try:
            print(json.dumps(resp.json(), indent=2))
        except Exception:
            print(resp.text)
    resp.raise_for_status()
    return resp.json() if resp.content else None


# ---------------------------
# AUTH
# ---------------------------

log("REGISTER")
email = f"test-{uuid.uuid4()}@example.com"

resp = session.post(
    f"{BASE_URL}/register",
    json={
        "firstname": "Test",
        "lastname": "User",
        "email": email,
        "password": "password123"
    }
)
check(resp)

log("LOGIN")
resp = session.post(
    f"{BASE_URL}/login",
    json={
        "email": email,
        "password": "password123"
    }
)
check(resp)

print("Cookies:", session.cookies.get_dict())


# ---------------------------
# ROOT
# ---------------------------

log("ROOT")
check(session.get(f"{BASE_URL}/"))


# ---------------------------
# COURSES
# ---------------------------

log("LIST COURSES")
courses = check(session.get(f"{BASE_URL}/courses"))

log("CREATE COURSE")
course = check(session.post(
    f"{BASE_URL}/courses",
    json={
        "name": "API Test Course",
        "description": "Created by automated test"
    }
))
course_id = course["uuid"]

log("GET COURSE DETAIL")
check(session.get(f"{BASE_URL}/courses/{course_id}"))

log("UPDATE COURSE")
check(session.put(
    f"{BASE_URL}/courses/{course_id}",
    json={"description": "Updated description", "name": "API TEST COURSE updated name!"}
))


# ---------------------------
# MATERIALS (URL material only)
# ---------------------------

log("LIST MATERIALS")
materials = check(session.get(
    f"{BASE_URL}/courses/{course_id}/materials"
))

log("CREATE URL MATERIAL")
material = check(session.post(
    f"{BASE_URL}/courses/{course_id}/materials",
    json={
        "type": "url",
        "name": "Official Docs",
        "description": "Project documentation",
        "url": "https://what.com/the/fucking/fuck/what/the/fuc/what/the/fuck/huh.shit"
    }
))
material_id = material["uuid"]

log("UPDATE URL MATERIAL")
check(session.put(
    f"{BASE_URL}/courses/{course_id}/materials/{material_id}",
    json={
        "name": "Updated Docs", "description": "updated descprition yeh", "url" : "www.example.com"
    }
))


log("CREATE FILE MATERIAL")
session.headers.pop("Content-Type", None)

file_path = "test.txt"
with open(file_path, "w") as f:
    f.write("Hello from API test")

file_mat = check(session.post(
    f"{BASE_URL}/courses/{course_id}/materials",
    files={
        "file": ("test.txt", open(file_path, "rb"), "text/plain"),
    },
    data={
        "type": "file",
        "name": "Test File Material",
        "description": "Uploaded by automated test",
    },
))

file_mat_uuid = file_mat["uuid"]


log("UPDATE FILE MATERIAL")
with open(file_path, "w") as f:
    f.write("Hello from API test UPDATE")

check(session.put(
    f"{BASE_URL}/courses/{course_id}/materials/{file_mat_uuid}",
    files={
        "file": ("test.txt", open(file_path, "rb"), "text/plain"),
    },
    data={
        "type": "file",
        "name": "Test File Material UPDATED",
        "description": "Uploaded by automated test AND UPDATED!",
    },
))
session.headers.update({"Content-Type": "application/json"})



log("DELETE MATERIAL")
check(session.delete(
    f"{BASE_URL}/courses/{course_id}/materials/{file_mat_uuid}"
))

log("DELETE MATERIAL")
check(session.delete(
    f"{BASE_URL}/courses/{course_id}/materials/{material_id}"
))


log("DELETE COURSE")
check(session.delete(f"{BASE_URL}/courses/{course_id}"))



quit()

# ---------------------------
# QUIZZES (NO SUBMIT)
# ---------------------------

log("LIST QUIZZES")
quizzes = check(session.get(
    f"{BASE_URL}/courses/{course_id}/quizzes"
))

log("CREATE QUIZ")
quiz = check(session.post(
    f"{BASE_URL}/courses/{course_id}/quizzes",
    json={
        "uuid": str(uuid.uuid4()),
        "title": "Test Quiz",
        "questions": [
            {
                "uuid": str(uuid.uuid4()),
                "type": "singleChoice",
                "question": "2 + 2?",
                "options": ["3", "4", "5"],
                "correctIndex": 1
            }
        ]
    }
))
quiz_id = quiz["uuid"]

log("GET QUIZ")
check(session.get(
    f"{BASE_URL}/courses/{course_id}/quizzes/{quiz_id}"
))

log("UPDATE QUIZ")
check(session.put(
    f"{BASE_URL}/courses/{course_id}/quizzes/{quiz_id}",
    json=quiz
))

log("DELETE QUIZ")
check(session.delete(
    f"{BASE_URL}/courses/{course_id}/quizzes/{quiz_id}"
))


# ---------------------------
# FEED (NO STREAM)
# ---------------------------

log("GET FEED")
check(session.get(
    f"{BASE_URL}/courses/{course_id}/feed"
))

log("CREATE FEED POST")
post = check(session.post(
    f"{BASE_URL}/courses/{course_id}/feed",
    json={"message": "Hello students"}
))
post_id = post["uuid"]

log("UPDATE FEED POST")
check(session.put(
    f"{BASE_URL}/courses/{course_id}/feed/{post_id}",
    json={"message": "Updated message", "edited": True}
))

log("DELETE FEED POST")
check(session.delete(
    f"{BASE_URL}/courses/{course_id}/feed/{post_id}"
))


# ---------------------------
# CLEANUP
# ---------------------------

log("DELETE COURSE")
check(session.delete(f"{BASE_URL}/courses/{course_id}"))

print("\n✅ ALL TESTS COMPLETED")
