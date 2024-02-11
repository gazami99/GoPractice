package main

import (
    "fmt"
    "runtime"
    "math/rand"
    "time"
    "sync"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func makeString(length int, charset string) string { // 랜덤문장배열
    b := make([]byte,length)
    for i := range b {
        b[i] = charset[seededRand.Intn(len(charset))]
    }

    return string(b)

}

func countChar(search_char string, strings string) int {  //  특정문장 찾는 함수 생성
    //
    count := 0

    for _ , i := range strings {

        if search_char == string(i) {
            count += 1
        }
    }

    return count
}

func main() {

    //Thread_num := runtime.NumCPU() // 기본적 Core 할당만큼 추가 그이상 추가시 성능상승 없음 
    Thread_num := 4             // 이유 외부 i/o  bound 가 아닐경우 parallelism 이 중요 
                                    // 쓰레드가 그이상 늘어날경우 Concurrency가  필요한경우  즉 외부 i/o 에 의존

    runtime.GOMAXPROCS(Thread_num)

    var wg sync.WaitGroup

    search_char := "A"         //검색키워드
    string_len := 1200000000  // 문장길이

    multi := int(string_len/Thread_num) // 분할배수  전체길이 / 쓰레드개수
    count_sum :=0                       // 분할누적합

    strings := makeString(string_len, charset)
    
    start := time.Now()

    num := countChar(search_char,strings) // 단순계산 

    end := time.Since(start)
    
    fmt.Println(end,num)

    //start gorutine

    channel := make(chan int)
    start = time.Now()

    for i :=1 ;  i <= Thread_num  ; i++ {

        wg.Add(1)

        go func(idx int) {
            start_num := multi*(idx-1)
            end_num := multi*idx
            channel <- countChar(search_char,strings[start_num:end_num])   // n개 분할 문장에대해 고루틴  생성
        }(i)

    }
    go func() { // A
    
		for count := range channel {
			count_sum += count
		    wg.Done()
		}
    
	}()


    wg.Wait()
    close(channel)

    end = time.Since(start)
    fmt.Println(end,count_sum)  // 두배 밖에 차이안나는지 생각해볼것

}
