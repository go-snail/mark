go get ":github.com/go-snail/mark"


//在main中调用Run方法，传入es所需值。

if err := mark.Run(&mark.ESConfig{Url:"",Scheme:"",Index:""},);err != nil {
        return err
    }

//在记日志处调用mark方法

mark.Mark(mark.Feilds{
        "tid":100004,
        "pid":123,
        "did":"12324qwrwer",
        "message":"error",
        "user":"lvxx",
        "type":1
    })