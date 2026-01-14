import unittest
import requests
import time
import re
from sseclient import SSEClient

BASE_URL = "http://localhost:80/api"  # <-- set your API base URL here

class TestPhase4(unittest.TestCase):
    test_course_id = None
    test_manual_post_id = None

    def test_01_create_course(self):
        payload = {
            "name": "Test Course for Phase 4",
            "description": "Testing course feed and SSE"
        }
        response = requests.post(f"{BASE_URL}/courses", json=payload)
        self.assertTrue(response.ok)
        data = response.json()
        self.assertIn("uuid", data)
        TestPhase4.test_course_id = data["uuid"]

    def test_02_list_feed_items_empty(self):
        self.assertIsNotNone(TestPhase4.test_course_id)
        response = requests.get(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed")
        self.assertTrue(response.ok)
        data = response.json()
        self.assertIsInstance(data, list)
        self.assertEqual(len(data), 0)

    def test_03_create_manual_post(self):
        self.assertIsNotNone(TestPhase4.test_course_id)
        new_post = {"message": "Welcome to the course! New materials will be published next week."}
        response = requests.post(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed", json=new_post)
        self.assertTrue(response.ok)
        data = response.json()
        self.assertEqual(data["type"], "manual")
        self.assertEqual(data["message"], new_post["message"])
        self.assertFalse(data.get("edited", False))
        self.assertIn("createdAt", data)
        TestPhase4.test_manual_post_id = data["uuid"]

    def test_04_retrieve_manual_post(self):
        self.assertIsNotNone(TestPhase4.test_course_id)
        self.assertIsNotNone(TestPhase4.test_manual_post_id)
        response = requests.get(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed")
        self.assertTrue(response.ok)
        data = response.json()
        post = next((item for item in data if item["uuid"] == TestPhase4.test_manual_post_id), None)
        self.assertIsNotNone(post)
        self.assertEqual(post["type"], "manual")
        self.assertEqual(post["message"], "Welcome to the course! New materials will be published next week.")

    def test_05_update_manual_post(self):
        self.assertIsNotNone(TestPhase4.test_course_id)
        self.assertIsNotNone(TestPhase4.test_manual_post_id)
        updated_post = {"message": "Updated: Materials will be published this Friday!", "edited": True}
        response = requests.put(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed/{TestPhase4.test_manual_post_id}", json=updated_post)
        self.assertTrue(response.ok)
        data = response.json()
        self.assertEqual(data["uuid"], TestPhase4.test_manual_post_id)
        self.assertEqual(data["message"], updated_post["message"])
        self.assertTrue(data["edited"])
        self.assertIn("updatedAt", data)

    def test_06_create_multiple_posts(self):
        self.assertIsNotNone(TestPhase4.test_course_id)
        posts = [
            {"message": "Quiz will be available next Monday"},
            {"message": "Office hours scheduled for Wednesday"},
            {"message": "Important: Assignment deadline extended"}
        ]
        for post in posts:
            response = requests.post(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed", json=post)
            self.assertTrue(response.ok)
            data = response.json()
            self.assertEqual(data["message"], post["message"])
            self.assertEqual(data["type"], "manual")
        feed_response = requests.get(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed")
        feed_data = feed_response.json()
        self.assertGreaterEqual(len(feed_data), 4)

    def test_07_delete_manual_post(self):
        self.assertIsNotNone(TestPhase4.test_course_id)
        self.assertIsNotNone(TestPhase4.test_manual_post_id)
        response = requests.delete(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed/{TestPhase4.test_manual_post_id}")
        self.assertEqual(response.status_code, 204)
        feed_response = requests.get(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed")
        feed_data = feed_response.json()
        deleted_post = next((item for item in feed_data if item["uuid"] == TestPhase4.test_manual_post_id), None)
        self.assertIsNone(deleted_post)

    def test_08_create_material_triggers_system_event(self):
        self.assertIsNotNone(TestPhase4.test_course_id)
        feed_before = requests.get(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed").json()
        before_count = len(feed_before)
        material = {
            "type": "url",
            "name": "Course Materials",
            "description": "Link to course materials",
            "url": "https://example.com/materials"
        }
        response = requests.post(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/materials", json=material)
        self.assertTrue(response.ok)
        feed_after = requests.get(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed").json()
        self.assertGreater(len(feed_after), before_count)
        system_event = next((item for item in feed_after if item["type"] == "system" and "material" in item["message"].lower()), None)
        self.assertIsNotNone(system_event)
        self.assertIn("createdAt", system_event)

    def test_09_create_quiz_triggers_system_event(self):
        self.assertIsNotNone(TestPhase4.test_course_id)
        feed_before = requests.get(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed").json()
        before_count = len(feed_before)
        quiz = {
            "title": "Test Quiz",
            "questions": [
                {"type": "singleChoice", "question": "What is 1 + 1?", "options": ["1", "2", "3", "4"], "correctIndex": 1}
            ]
        }
        response = requests.post(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/quizzes", json=quiz)
        self.assertTrue(response.ok)
        feed_after = requests.get(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed").json()
        self.assertGreater(len(feed_after), before_count)
        system_event = next((item for item in feed_after if item["type"] == "system" and "quiz" in item["message"].lower()), None)
        self.assertIsNotNone(system_event)
        self.assertIn("createdAt", system_event)

    def test_10_course_detail_includes_feed(self):
        self.assertIsNotNone(TestPhase4.test_course_id)
        response = requests.get(f"{BASE_URL}/courses/{TestPhase4.test_course_id}")
        self.assertTrue(response.ok)
        data = response.json()
        self.assertIn("feed", data)
        self.assertIsInstance(data["feed"], list)
        self.assertGreaterEqual(len(data["feed"]), 1)
        for item in data["feed"]:
            self.assertIn("uuid", item)
            self.assertIn(item["type"], ["manual", "system"])
            self.assertIn("message", item)
            self.assertIn("createdAt", item)

    def test_11_uuid_format(self):
        self.assertIsNotNone(TestPhase4.test_course_id)
        response = requests.post(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed", json={"message": "UUID test post"})
        self.assertTrue(response.ok)
        data = response.json()
        uuid_pattern = re.compile(r"^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$", re.I)
        self.assertTrue(uuid_pattern.match(data["uuid"]))

    def test_12_valid_timestamps(self):
        self.assertIsNotNone(TestPhase4.test_course_id)
        create_data = requests.post(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed", json={"message": "Timestamp test post"}).json()
        self.assertIsNotNone(create_data["createdAt"])
        time.sleep(0.1)
        update_data = requests.put(
            f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed/{create_data['uuid']}",
            json={"message": "Updated timestamp test post", "edited": True}
        ).json()
        self.assertIsNotNone(update_data["updatedAt"])

    def test_13_sse_connection(self):
        self.assertIsNotNone(TestPhase4.test_course_id)
        stream_url = f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed/stream"
        # Just check connection headers
        response = requests.get(stream_url, headers={"Accept": "text/event-stream"}, stream=True)
        self.assertTrue(response.ok)
        self.assertIn("text/event-stream", response.headers.get("content-type", ""))
        response.close()

    def test_14_feed_ordering(self):
        self.assertIsNotNone(TestPhase4.test_course_id)
        requests.post(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed", json={"message": "First post"})
        time.sleep(1)
        requests.post(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed", json={"message": "Second post"})
        time.sleep(1)
        requests.post(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed", json={"message": "Third post"})
        feed_data = requests.get(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed").json()
        third_index = next(i for i, item in enumerate(feed_data) if item["message"] == "Third post")
        second_index = next(i for i, item in enumerate(feed_data) if item["message"] == "Second post")
        first_index = next(i for i, item in enumerate(feed_data) if item["message"] == "First post")
        self.assertLess(third_index, second_index)
        self.assertLess(second_index, first_index)

    def test_15_sse_receives_new_items(self):
        self.assertIsNotNone(TestPhase4.test_course_id)
        stream_url = f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed/stream"
        received_event = False

        def listen_sse():
            nonlocal received_event
            for event in SSEClient(stream_url, timeout=10):
                if event.data:
                    received_event = True
                    break

        import threading
        listener = threading.Thread(target=listen_sse)
        listener.start()
        time.sleep(0.5)
        requests.post(f"{BASE_URL}/courses/{TestPhase4.test_course_id}/feed", json={"message": "SSE test post"})
        listener.join(timeout=10)
        self.assertTrue(received_event)

    def test_16_cleanup_course(self):
        self.assertIsNotNone(TestPhase4.test_course_id)
        response = requests.delete(f"{BASE_URL}/courses/{TestPhase4.test_course_id}")
        self.assertEqual(response.status_code, 204)


if __name__ == "__main__":
    unittest.main()
