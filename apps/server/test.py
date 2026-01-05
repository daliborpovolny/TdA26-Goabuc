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


import time
import mimetypes
from pathlib import Path

log("PHASE 2 – MATERIALS")

# --------------------------
# CREATE TEST COURSE
# ---------------------------

course = check(session.post(
    f"{BASE_URL}/courses",
    json={
        "name": "Test Course for Phase 2",
        "description": "Testing course materials"
    }
))
course_id = course["uuid"]


# ---------------------------
# LIST MATERIALS (EMPTY)
# ---------------------------

log("LIST MATERIALS (EMPTY)")
materials = check(session.get(
    f"{BASE_URL}/courses/{course_id}/materials"
))
assert materials == []


# ---------------------------
# ADD URL MATERIAL
# ---------------------------

log("ADD URL MATERIAL")
url_material = check(session.post(
    f"{BASE_URL}/courses/{course_id}/materials",
    json={
        "type": "url",
        "name": "Official Documentation",
        "description": "Link to docs",
        "url": "https://example.com/docs"
    }
))
url_material_id = url_material["uuid"]


# ---------------------------
# ADD FILE MATERIAL (PDF)
# ---------------------------

log("ADD FILE MATERIAL – PDF")
session.headers.pop("Content-Type", None)

pdf_path = Path("test.pdf")

with open(pdf_path, "rb") as f:
    file_material = check(session.post(
        f"{BASE_URL}/courses/{course_id}/materials",
        files={"file": ("test.pdf", f, "application/pdf")},
        data={
            "type": "file",
            "name": "Course Syllabus",
            "description": "PDF syllabus"
        }
    ))

file_material_id = file_material["uuid"]


# ---------------------------
# ADD FILE MATERIAL (IMAGE)
# ---------------------------

log("ADD FILE MATERIAL – IMAGE")
img_path = Path("test.png")

with open(img_path, "rb") as f:
    check(session.post(
        f"{BASE_URL}/courses/{course_id}/materials",
        files={"file": ("test.png", f, "image/png")},
        data={
            "type": "file",
            "name": "Diagram",
            "description": "PNG diagram"
        }
    ))


# ---------------------------
# REJECT LARGE FILE (>30MB)
# ---------------------------

log("REJECT LARGE FILE")
big_path = Path("big.bin")
big_path.write_bytes(b"\x00" * (31 * 1024 * 1024))

with open(big_path, "rb") as f:
    resp = session.post(
        f"{BASE_URL}/courses/{course_id}/materials",
        files={"file": ("big.bin", f, "application/octet-stream")},
        data={
            "type": "file",
            "name": "Too Large",
            "description": "Should fail"
        }
    )
    assert resp.status_code == 400


# ---------------------------
# REJECT UNSUPPORTED TYPE
# ---------------------------

log("REJECT UNSUPPORTED FILE TYPE")
with open("malware.exe", "wb") as f:
    f.write(b"MZ")

with open("malware.exe", "rb") as f:
    resp = session.post(
        f"{BASE_URL}/courses/{course_id}/materials",
        files={"file": ("malware.exe", f, "application/x-msdownload")},
        data={
            "type": "file",
            "name": "Bad File",
            "description": "Should fail"
        }
    )
    assert resp.status_code == 400

session.headers.update({"Content-Type": "application/json"})

# ---------------------------
# LIST MATERIALS (COUNT + ORDER)
# ---------------------------

log("LIST MATERIALS (ORDER)")
check(session.post(
    f"{BASE_URL}/courses/{course_id}/materials",
    json={
        "type": "url",
        "name": "First",
        "description": "First",
        "url": "https://example.com/1"
    }
))
time.sleep(1)

check(session.post(
    f"{BASE_URL}/courses/{course_id}/materials",
    json={
        "type": "url",
        "name": "Second",
        "description": "Second",
        "url": "https://example.com/2"
    }
))
time.sleep(1)

check(session.post(
    f"{BASE_URL}/courses/{course_id}/materials",
    json={
        "type": "url",
        "name": "Third",
        "description": "Third",
        "url": "https://example.com/3"
    }
))

materials = check(session.get(
    f"{BASE_URL}/courses/{course_id}/materials"
))

names = [m["name"] for m in materials]
assert names.index("Third") < names.index("Second") < names.index("First")


# ---------------------------
# COURSE DETAIL INCLUDES MATERIALS
# ---------------------------

log("COURSE DETAIL INCLUDES MATERIALS")
course_detail = check(session.get(
    f"{BASE_URL}/courses/{course_id}"
))
assert isinstance(course_detail["materials"], list)
assert len(course_detail["materials"]) >= 3


# ---------------------------
# UPDATE URL MATERIAL
# ---------------------------

log("UPDATE URL MATERIAL")
check(session.put(
    f"{BASE_URL}/courses/{course_id}/materials/{url_material_id}",
    json={
        "name": "Updated Docs",
        "description": "Updated",
        "url": "https://example.com/new"
    }
))


# ---------------------------
# UPDATE FILE MATERIAL METADATA
# ---------------------------

session.headers.pop("Content-Type", None)


log("UPDATE FILE MATERIAL METADATA")
check(session.put(
    f"{BASE_URL}/courses/{course_id}/materials/{file_material_id}",
    json={
        "name": "Updated Syllabus",
        "description": "Updated description"
    }
))


# ---------------------------
# REPLACE FILE
# ---------------------------

log("REPLACE FILE")
with open(pdf_path, "rb") as f:
    check(session.put(
        f"{BASE_URL}/courses/{course_id}/materials/{file_material_id}",
        files={"file": ("updated.pdf", f, "application/pdf")},
        data={"name": "Replaced Syllabus"}
    ))

session.headers.update({"Content-Type": "application/json"})


# ---------------------------
# DELETE MATERIALS
# ---------------------------

log("DELETE MATERIALS")
session.delete(
    f"{BASE_URL}/courses/{course_id}/materials/{url_material_id}"
)
session.delete(
    f"{BASE_URL}/courses/{course_id}/materials/{file_material_id}"
)



# ---------------------------
# ACCEPT ALL SUPPORTED FILE FORMATS
# ---------------------------

log("ACCEPT ALL SUPPORTED FILE FORMATS")
session.headers.pop("Content-Type", None)

supported_formats = [
    {
        "ext": "txt",
        "mime": "text/plain",
        "content": b"Text content",
    },
    {
        "ext": "jpg",
        "mime": "image/jpeg",
        "content": bytes([0xFF, 0xD8, 0xFF]),
    },
    {
        "ext": "jpeg",
        "mime": "image/jpeg",
        "content": bytes([0xFF, 0xD8, 0xFF]),
    },
    {
        "ext": "gif",
        "mime": "image/gif",
        "content": b"GIF89a",
    },
    {
        "ext": "mp3",
        "mime": "audio/mpeg",
        "content": b"ID3",
    },
    {
        "ext": "mp4",
        "mime": "video/mp4",
        "content": bytes([0x00, 0x00, 0x00, 0x18]),
    },
    {
        "ext": "docx",
        "mime": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
        "content": b"PK",
    },
]

for fmt in supported_formats:
    filename = f"test.{fmt['ext']}"

    with open(filename, "wb") as f:
        f.write(fmt["content"])

    with open(filename, "rb") as f:
        resp = session.post(
            f"{BASE_URL}/courses/{course_id}/materials",
            files={
                "file": (filename, f, fmt["mime"]),
            },
            data={
                "type": "file",
                "name": f"Test {fmt['ext'].upper()} File",
                "description": f"Testing {fmt['ext']} format",
            },
        )

    print(f"→ {filename}: {resp.status_code}")
    data = check(resp)

    assert data["type"] == "file"
    assert data["mimeType"] == fmt["mime"]



# ---------------------------
# CLEANUP COURSE
# ---------------------------

log("DELETE COURSE")
session.delete(f"{BASE_URL}/courses/{course_id}")





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
