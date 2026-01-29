#!/bin/bash

API="http://localhost:8080/submit-answer"
EXAM_ID="exam170126"

post_answer () {
  curl -s -X POST "$API" \
    -H "Content-Type: application/json" \
    -d "{
      \"examId\": \"$EXAM_ID\",
      \"questionId\": \"$1\",
      \"ans\": \"$2\",
      \"StudentID\": \"$3\"
    }" > /dev/null
}

echo "=== Starting normal MCQ submissions ==="

# ---------- NORMAL STUDENTS ----------
NORMAL_STUDENTS=(s1 s2 s4 s5 s6 s8 s9 s10)
OPTIONS=(A B C D)

random_option () {
  echo "${OPTIONS[$RANDOM % 4]}"
}

for student in "${NORMAL_STUDENTS[@]}"; do
  post_answer Q1 "Option $(random_option)" "$student"
  sleep 1

  post_answer Q2 "Option $(random_option)" "$student"
  sleep 1

  post_answer Q3 "Option $(random_option)" "$student"
  sleep 1
done

echo "=== Simulating MCQ copying behavior ==="

# ---------- CHEATING STUDENTS ----------
# s3 and s7 copy ONLY Q2

# First attempt (both unsure, pick same wrong option)
post_answer Q2 "Option B" "s3"
sleep 1
post_answer Q2 "Option B" "s7"

# Second edit (both switch to correct option almost together)
sleep 2
post_answer Q2 "Option C" "s3"
sleep 1
post_answer Q2 "Option C" "s7"

# Legit answers for other questions
post_answer Q1 "Option A" "s3"
sleep 1
post_answer Q3 "Option D" "s3"

post_answer Q1 "Option B" "s7"
sleep 1
post_answer Q3 "Option A" "s7"

echo "=== MCQ Simulation complete ==="
